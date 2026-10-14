package blog

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CryptoCommunityBlogMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls" json:"screenshot_urls"`
	VideoURL       string   `bson:"video_url" json:"video_url"`
}

type CryptoCommunityBlogAuthor struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	ImageURL string `bson:"image_url" json:"image_url"`
}

type CryptoCommunityBlogSubsection struct {
	ID         string   `bson:"_id" json:"_id"`
	Name       string   `bson:"name" json:"name"`
	Paragraphs []string `bson:"paragraphs" json:"paragraphs"`
}

type CryptoCommunityBlogScreenshot struct {
	Index       int    `bson:"index" json:"index"`
	Description string `bson:"description" json:"description"`
}

type CryptoCommunityBlogVideo struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type CryptoCommunityBlogSection struct {
	ID          string                          `bson:"_id" json:"_id"`
	Name        string                          `bson:"name" json:"name"`
	Paragraphs  []string                        `bson:"paragraphs" json:"paragraphs"`
	Subsections []CryptoCommunityBlogSubsection `bson:"subsections,omitempty" json:"subsections,omitempty"`
	Screenshot  *CryptoCommunityBlogScreenshot  `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	Video       *CryptoCommunityBlogVideo       `bson:"video,omitempty" json:"video,omitempty"`
}

type CryptoCommunityBlog struct {
	ID          primitive.ObjectID           `bson:"_id,omitempty" json:"_id,omitempty"`
	Slug        string                       `bson:"slug" json:"slug"`
	Name        string                       `bson:"name" json:"name"`
	Description string                       `bson:"description" json:"description"`
	ImageURL    string                       `bson:"image_url" json:"image_url"`
	Categories  []string                     `bson:"categories" json:"categories"`
	Media       CryptoCommunityBlogMedia     `bson:"media" json:"media"`
	AddedBy     CryptoCommunityBlogAuthor    `bson:"added_by" json:"added_by"`
	CreatedAt   time.Time                    `bson:"created_at" json:"created_at"`
	Content     []CryptoCommunityBlogSection `bson:"content" json:"content"`
}

type cryptoCommunityBlogRequest struct {
	CryptoCommunityBlog
	ID json.RawMessage `json:"_id"`
}
