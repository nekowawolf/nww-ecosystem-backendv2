package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatorsBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type CreatorsBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type CreatorsBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type CreatorsBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type CreatorsBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type CreatorsBlogSection struct {
	ID          string                   `bson:"_id" json:"_id"`
	Name        string                   `bson:"name" json:"name"`
	Paragraphs  []string                 `bson:"paragraphs" json:"paragraphs"`
	Subsections []CreatorsBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *CreatorsBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *CreatorsBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type CreatorsBlog struct {
	ID          primitive.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string                `bson:"slug" json:"slug"`
	Name        string                `bson:"name" json:"name"`
	Description string                `bson:"description" json:"description"`
	ImageURL    string                `bson:"image_url" json:"image_url"`
	Categories  []string              `bson:"categories" json:"categories"`
	Media       CreatorsBlogMedia     `bson:"media" json:"media"`
	AddedBy     CreatorsBlogAuthor    `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time             `bson:"created_at" json:"created_at"`
	Content     []CreatorsBlogSection `bson:"content" json:"content"`
}

type creatorsBlogRequest struct {
	CreatorsBlog
	ID json.RawMessage `json:"_id"`
}
