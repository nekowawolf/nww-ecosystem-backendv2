package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Web3Tools struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	Chains      []string           `bson:"chains,omitempty" json:"chains,omitempty"`
	ImageURL    string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Twitter     string             `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Instagram   string             `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Discord     string             `bson:"discord,omitempty" json:"discord,omitempty"`
	Telegram    string             `bson:"telegram,omitempty" json:"telegram,omitempty"`
	Youtube     string             `bson:"youtube,omitempty" json:"youtube,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}