package blog

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const githubBlogCollection = "githubBlogs"

var githubBlogSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func normalizeGithubBlog(blog *GithubBlog) {
	if blog.Categories == nil {
		blog.Categories = []string{}
	}
	if blog.Media.ScreenshotURLs == nil {
		blog.Media.ScreenshotURLs = []string{}
	}
	if blog.Content == nil {
		blog.Content = []GithubBlogSection{}
	}
	for i := range blog.Content {
		section := &blog.Content[i]
		if section.Paragraphs == nil {
			section.Paragraphs = []string{}
		}
		for j := range section.Subsections {
			if section.Subsections[j].Paragraphs == nil {
				section.Subsections[j].Paragraphs = []string{}
			}
		}
	}
}

func validateGithubBlog(blog *GithubBlog) error {
	blog.Slug = strings.TrimSpace(blog.Slug)
	blog.Name = strings.TrimSpace(blog.Name)
	normalizeGithubBlog(blog)
	if len(blog.Slug) > 200 || !githubBlogSlugPattern.MatchString(blog.Slug) {
		return fmt.Errorf("slug must use lowercase letters, numbers, and single hyphens (maximum 200 characters)")
	}
	if blog.Name == "" || strings.TrimSpace(blog.Description) == "" || strings.TrimSpace(blog.ImageURL) == "" {
		return fmt.Errorf("name, description, and image_url are required")
	}
	if strings.TrimSpace(blog.AddedBy.Name) == "" {
		return fmt.Errorf("added_by.name is required")
	}
	if len(blog.Content) == 0 {
		return fmt.Errorf("content must contain at least one section")
	}
	ids := make(map[string]bool)
	validateHeading := func(id, name string) error {
		if strings.TrimSpace(id) == "" || strings.ContainsAny(id, " \t\r\n") || strings.TrimSpace(name) == "" {
			return fmt.Errorf("each section and subsection requires a non-whitespace _id and a name")
		}
		if ids[id] {
			return fmt.Errorf("duplicate section or subsection _id: %s", id)
		}
		ids[id] = true
		return nil
	}
	for _, section := range blog.Content {
		if err := validateHeading(section.ID, section.Name); err != nil {
			return err
		}
		if section.Screenshot != nil && (section.Screenshot.Index < 0 || section.Screenshot.Index >= len(blog.Media.ScreenshotURLs)) {
			return fmt.Errorf("screenshot index in section %s is outside media.screenshot_urls", section.ID)
		}
		if section.Video != nil && (strings.TrimSpace(blog.Media.VideoURL) == "" || strings.TrimSpace(section.Video.Name) == "") {
			return fmt.Errorf("video in section %s requires a name and media.video_url", section.ID)
		}
		for _, subsection := range section.Subsections {
			if err := validateHeading(subsection.ID, subsection.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func InsertGithubBlog(blog GithubBlog) (*GithubBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	blog.ID = primitive.NewObjectID()
	if blog.CreatedAt.IsZero() {
		blog.CreatedAt = time.Now().UTC()
	}
	normalizeGithubBlog(&blog)
	if _, err := config.Database.Collection(githubBlogCollection).InsertOne(ctx, blog); err != nil {
		return nil, fmt.Errorf("error inserting github blog: %w", err)
	}
	return &blog, nil
}

func GetAllGithubBlogs() ([]GithubBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: -1}})
	cursor, err := config.Database.Collection(githubBlogCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving github blogs: %w", err)
	}
	defer cursor.Close(ctx)

	blogs := []GithubBlog{}
	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, fmt.Errorf("error decoding github blogs: %w", err)
	}
	for i := range blogs {
		normalizeGithubBlog(&blogs[i])
	}
	return blogs, nil
}

func getGithubBlog(filter bson.M) (*GithubBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	var blog GithubBlog
	if err := config.Database.Collection(githubBlogCollection).FindOne(ctx, filter).Decode(&blog); err != nil {
		return nil, err
	}
	normalizeGithubBlog(&blog)
	return &blog, nil
}

func GetGithubBlogByID(id primitive.ObjectID) (*GithubBlog, error) {
	return getGithubBlog(bson.M{"_id": id})
}

func GetGithubBlogBySlug(slug string) (*GithubBlog, error) {
	return getGithubBlog(bson.M{"slug": slug})
}

func UpdateGithubBlogByID(id primitive.ObjectID, blog GithubBlog) (*GithubBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	normalizeGithubBlog(&blog)
	update := bson.M{"$set": bson.M{
		"slug": blog.Slug, "name": blog.Name, "description": blog.Description,
		"image_url": blog.ImageURL, "categories": blog.Categories,
		"media": blog.Media, "added_by": blog.AddedBy, "content": blog.Content,
	}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated GithubBlog
	if err := config.Database.Collection(githubBlogCollection).FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&updated); err != nil {
		return nil, err
	}
	normalizeGithubBlog(&updated)
	return &updated, nil
}

func DeleteGithubBlogByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	result, err := config.Database.Collection(githubBlogCollection).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}