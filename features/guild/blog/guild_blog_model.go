package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GuildBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type GuildBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type GuildBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type GuildBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type GuildBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type GuildBlogSection struct {
	ID          string                `bson:"_id" json:"_id"`
	Name        string                `bson:"name" json:"name"`
	Paragraphs  []string              `bson:"paragraphs" json:"paragraphs"`
	Subsections []GuildBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *GuildBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *GuildBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type GuildBlog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string             `bson:"slug" json:"slug"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	Categories  []string           `bson:"categories" json:"categories"`
	Media       GuildBlogMedia     `bson:"media" json:"media"`
	AddedBy     GuildBlogAuthor    `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	Content     []GuildBlogSection `bson:"content" json:"content"`
}

type guildBlogRequest struct {
	GuildBlog
	ID json.RawMessage `json:"_id"`
}
