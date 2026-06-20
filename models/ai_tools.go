package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AITools struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Categories  []string           `bson:"categories,omitempty" json:"categories,omitempty"`
	VideoURL    string             `bson:"video_url,omitempty" json:"video_url,omitempty"`
	ImgURL      string             `bson:"imgURL,omitempty" json:"imgURL,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Twitter     string             `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Instagram   string             `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Discord     string             `bson:"discord,omitempty" json:"discord,omitempty"`
	Telegram    string             `bson:"telegram,omitempty" json:"telegram,omitempty"`
}