package web3_tools

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Web3ToolsMedia struct {
	ScreenshotURLs []string `bson:"screenshot_urls,omitempty" json:"screenshot_urls,omitempty"`
	VideoURL       string   `bson:"video_url,omitempty" json:"video_url,omitempty"`
}

type AddedByInfo struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	URL  string `bson:"url,omitempty" json:"url,omitempty"`
}

type Web3Tools struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	Chains      []string           `bson:"chains,omitempty" json:"chains,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Media       Web3ToolsMedia     `bson:"media,omitempty" json:"media,omitempty"`
	AddedBy     *AddedByInfo       `bson:"added_by,omitempty" json:"added_by,omitempty"`
	Twitter     string             `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Instagram   string             `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Discord     string             `bson:"discord,omitempty" json:"discord,omitempty"`
	Telegram    string             `bson:"telegram,omitempty" json:"telegram,omitempty"`
	Youtube     string             `bson:"youtube,omitempty" json:"youtube,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}