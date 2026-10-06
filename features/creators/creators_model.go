package creators

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AddedByInfo struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	URL  string `bson:"url,omitempty" json:"url,omitempty"`
}

type CreatorsSocials struct {
	Twitter       string `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Instagram     string `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Discord       string `bson:"discord,omitempty" json:"discord,omitempty"`
	Youtube       string `bson:"youtube,omitempty" json:"youtube,omitempty"`
	Telegram      string `bson:"telegram,omitempty" json:"telegram,omitempty"`
	Github        string `bson:"github,omitempty" json:"github,omitempty"`
	Tiktok        string `bson:"tiktok,omitempty" json:"tiktok,omitempty"`
}

type CreatorsPlatforms struct {
	Fiverr        string `bson:"fiverr,omitempty" json:"fiverr,omitempty"`
	Upwork        string `bson:"upwork,omitempty" json:"upwork,omitempty"`
	PeoplePerHour string `bson:"peopleperhour,omitempty" json:"peopleperhour,omitempty"`
	Freelancer    string `bson:"freelancer,omitempty" json:"freelancer,omitempty"`
}

type Creators struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	Language    string             `bson:"language,omitempty" json:"language,omitempty"`
	OpenToWork  bool               `bson:"open_to_work,omitempty" json:"open_to_work,omitempty"`
	Socials     CreatorsSocials    `bson:"socials,omitempty" json:"socials,omitempty"`
	Platforms   CreatorsPlatforms  `bson:"platforms,omitempty" json:"platforms,omitempty"`
	AddedBy     *AddedByInfo       `bson:"added_by,omitempty" json:"added_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}