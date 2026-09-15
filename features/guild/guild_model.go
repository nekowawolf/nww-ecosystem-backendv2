package guild

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Socials struct {
	Twitter   string `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Instagram string `bson:"instagram,omitempty" json:"instagram,omitempty"`
	Discord   string `bson:"discord,omitempty" json:"discord,omitempty"`
	Github    string `bson:"github,omitempty" json:"github,omitempty"`
	Youtube   string `bson:"youtube,omitempty" json:"youtube,omitempty"`
}

type Guild struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty" json:"image_url,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Platform    string             `bson:"platform,omitempty" json:"platform,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	Link        string             `bson:"link,omitempty" json:"link,omitempty"`
	Socials     Socials            `bson:"socials,omitempty" json:"socials,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}