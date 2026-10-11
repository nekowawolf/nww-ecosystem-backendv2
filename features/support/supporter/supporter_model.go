package supporter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Supporter struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name,omitempty" json:"name,omitempty"`
	URL       string             `bson:"url,omitempty" json:"url,omitempty"`
	Platform  string             `bson:"platform,omitempty" json:"platform,omitempty"`
	Proof     string             `bson:"proof,omitempty" json:"proof,omitempty"`
	Amount    int64              `bson:"amount,omitempty" json:"amount,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}