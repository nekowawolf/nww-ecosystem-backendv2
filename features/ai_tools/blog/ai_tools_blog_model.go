package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIToolsBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type AIToolsBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type AIToolsBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type AIToolsBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type AIToolsBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type AIToolsBlogSection struct {
	ID          string                  `bson:"_id" json:"_id"`
	Name        string                  `bson:"name" json:"name"`
	Paragraphs  []string                `bson:"paragraphs" json:"paragraphs"`
	Subsections []AIToolsBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *AIToolsBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *AIToolsBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type AIToolsBlog struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string               `bson:"slug" json:"slug"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	ImageURL    string               `bson:"image_url" json:"image_url"`
	Categories  []string             `bson:"categories" json:"categories"`
	Media       AIToolsBlogMedia     `bson:"media" json:"media"`
	AddedBy     AIToolsBlogAuthor    `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	Content     []AIToolsBlogSection `bson:"content" json:"content"`
}

type aiToolsBlogRequest struct {
	AIToolsBlog
	ID json.RawMessage `json:"_id"`
}