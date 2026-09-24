package community

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AddedByInfo struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	URL  string `bson:"url,omitempty" json:"url,omitempty"`
}

type Socials struct {
	Twitter   string `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Discord   string `bson:"discord,omitempty" json:"discord,omitempty"`
	Github    string `bson:"github,omitempty" json:"github,omitempty"`
	Youtube   string `bson:"youtube,omitempty" json:"youtube,omitempty"`
}

type CryptoCommunity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Platforms   string             `bson:"platforms,omitempty" json:"platforms,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Link        string             `bson:"link,omitempty" json:"link,omitempty"`
	Socials     Socials            `bson:"socials,omitempty" json:"socials,omitempty"`
	AddedBy     *AddedByInfo       `bson:"added_by,omitempty" json:"added_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}