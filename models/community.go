package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CryptoCommunity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	Platforms string             `bson:"platforms,omitempty" json:"platforms,omitempty"`
	Category  string             `bson:"category,omitempty" json:"category,omitempty"`
	ImgURL    string             `bson:"img_url,omitempty" json:"img_url,omitempty"`
	LinkURL   string             `bson:"link_url,omitempty" json:"link_url,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}