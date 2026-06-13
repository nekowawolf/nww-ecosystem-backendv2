package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Web3Tool struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	Chains      []string           `bson:"chains,omitempty" json:"chains,omitempty"`
	ImageURL    string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Twitter     string             `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Discord     string             `bson:"discord,omitempty" json:"discord,omitempty"`
	Telegram    string             `bson:"telegram,omitempty" json:"telegram,omitempty"`
}