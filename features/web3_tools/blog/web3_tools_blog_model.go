package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Web3ToolsBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type Web3ToolsBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type Web3ToolsBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type Web3ToolsBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type Web3ToolsBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type Web3ToolsBlogSection struct {
	ID          string                    `bson:"_id" json:"_id"`
	Name        string                    `bson:"name" json:"name"`
	Paragraphs  []string                  `bson:"paragraphs" json:"paragraphs"`
	Subsections []Web3ToolsBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *Web3ToolsBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *Web3ToolsBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type Web3ToolsBlog struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string                 `bson:"slug" json:"slug"`
	Name        string                 `bson:"name" json:"name"`
	Description string                 `bson:"description" json:"description"`
	ImageURL    string                 `bson:"image_url" json:"image_url"`
	Categories  []string               `bson:"categories" json:"categories"`
	Media       Web3ToolsBlogMedia     `bson:"media" json:"media"`
	AddedBy     Web3ToolsBlogAuthor    `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	Content     []Web3ToolsBlogSection `bson:"content" json:"content"`
}

type web3ToolsBlogRequest struct {
	Web3ToolsBlog
	ID json.RawMessage `json:"_id"`
}
