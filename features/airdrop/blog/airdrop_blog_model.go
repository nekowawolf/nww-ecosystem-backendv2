package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AirdropBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type AirdropBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type AirdropBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type AirdropBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type AirdropBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type AirdropBlogSection struct {
	ID          string                  `bson:"_id" json:"_id"`
	Name        string                  `bson:"name" json:"name"`
	Paragraphs  []string                `bson:"paragraphs" json:"paragraphs"`
	Subsections []AirdropBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *AirdropBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *AirdropBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type AirdropBlog struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string               `bson:"slug" json:"slug"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	ImageURL    string               `bson:"image_url" json:"image_url"`
	Categories  []string             `bson:"categories" json:"categories"`
	Media       AirdropBlogMedia     `bson:"media" json:"media"`
	AddedBy     AirdropBlogAuthor    `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	Content     []AirdropBlogSection `bson:"content" json:"content"`
}

type airdropBlogRequest struct {
	AirdropBlog
	ID json.RawMessage `json:"_id"`
}