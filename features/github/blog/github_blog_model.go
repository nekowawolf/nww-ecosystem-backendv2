package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GithubBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type GithubBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type GithubBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type GithubBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type GithubBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type GithubBlogSection struct {
	ID          string                 `bson:"_id" json:"_id"`
	Name        string                 `bson:"name" json:"name"`
	Paragraphs  []string               `bson:"paragraphs" json:"paragraphs"`
	Subsections []GithubBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *GithubBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *GithubBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type GithubBlog struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string              `bson:"slug" json:"slug"`
	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
	ImageURL    string              `bson:"image_url" json:"image_url"`
	Categories  []string            `bson:"categories" json:"categories"`
	Media       GithubBlogMedia     `bson:"media" json:"media"`
	AddedBy     GithubBlogAuthor    `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	Content     []GithubBlogSection `bson:"content" json:"content"`
}

type githubBlogRequest struct {
	GithubBlog
	ID json.RawMessage `json:"_id"`
}