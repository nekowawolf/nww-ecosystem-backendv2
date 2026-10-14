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

const airdropBlogCollection = "airdropBlogs"

var airdropBlogSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func normalizeAirdropBlog(blog *AirdropBlog) {
	if blog.Categories == nil {
		blog.Categories = []string{}
	}
	if blog.Media.ScreenshotURLs == nil {
		blog.Media.ScreenshotURLs = []string{}
	}
	if blog.Content == nil {
		blog.Content = []AirdropBlogSection{}
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

func validateAirdropBlog(blog *AirdropBlog) error {
	blog.Slug = strings.TrimSpace(blog.Slug)
	blog.Name = strings.TrimSpace(blog.Name)
	normalizeAirdropBlog(blog)
	if len(blog.Slug) > 200 || !airdropBlogSlugPattern.MatchString(blog.Slug) {
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

func InsertAirdropBlog(blog AirdropBlog) (*AirdropBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	blog.ID = primitive.NewObjectID()
	if blog.CreatedAt.IsZero() {
		blog.CreatedAt = time.Now().UTC()
	}
	normalizeAirdropBlog(&blog)
	if _, err := config.Database.Collection(airdropBlogCollection).InsertOne(ctx, blog); err != nil {
		return nil, fmt.Errorf("error inserting airdrop blog: %w", err)
	}
	return &blog, nil
}

func GetAllAirdropBlogs() ([]AirdropBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: -1}})
	cursor, err := config.Database.Collection(airdropBlogCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving airdrop blogs: %w", err)
	}
	defer cursor.Close(ctx)

	blogs := []AirdropBlog{}
	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, fmt.Errorf("error decoding airdrop blogs: %w", err)
	}
	for i := range blogs {
		normalizeAirdropBlog(&blogs[i])
	}
	return blogs, nil
}

func getAirdropBlog(filter bson.M) (*AirdropBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	var blog AirdropBlog
	if err := config.Database.Collection(airdropBlogCollection).FindOne(ctx, filter).Decode(&blog); err != nil {
		return nil, err
	}
	normalizeAirdropBlog(&blog)
	return &blog, nil
}

func GetAirdropBlogByID(id primitive.ObjectID) (*AirdropBlog, error) {
	return getAirdropBlog(bson.M{"_id": id})
}

func GetAirdropBlogBySlug(slug string) (*AirdropBlog, error) {
	return getAirdropBlog(bson.M{"slug": slug})
}

func UpdateAirdropBlogByID(id primitive.ObjectID, blog AirdropBlog) (*AirdropBlog, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	normalizeAirdropBlog(&blog)
	update := bson.M{"$set": bson.M{
		"slug": blog.Slug, "name": blog.Name, "description": blog.Description,
		"image_url": blog.ImageURL, "categories": blog.Categories,
		"media": blog.Media, "added_by": blog.AddedBy, "content": blog.Content,
	}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated AirdropBlog
	if err := config.Database.Collection(airdropBlogCollection).FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&updated); err != nil {
		return nil, err
	}
	normalizeAirdropBlog(&updated)
	return &updated, nil
}

func DeleteAirdropBlogByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	result, err := config.Database.Collection(airdropBlogCollection).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}