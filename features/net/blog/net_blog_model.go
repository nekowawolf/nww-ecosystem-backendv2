package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NetBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type NetBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type NetBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type NetBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type NetBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type NetBlogSection struct {
	ID          string              `bson:"_id" json:"_id"`
	Name        string              `bson:"name" json:"name"`
	Paragraphs  []string            `bson:"paragraphs" json:"paragraphs"`
	Subsections []NetBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *NetBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *NetBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type NetBlog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string             `bson:"slug" json:"slug"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	Categories  []string           `bson:"categories" json:"categories"`
	Media       NetBlogMedia       `bson:"media" json:"media"`
	AddedBy     NetBlogAuthor      `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	Content     []NetBlogSection   `bson:"content" json:"content"`
}

type netBlogRequest struct {
	NetBlog
	ID json.RawMessage `json:"_id"`
}
