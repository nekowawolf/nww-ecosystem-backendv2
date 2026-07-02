package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title     string             `bson:"title,omitempty" json:"title,omitempty"`
	Content   string             `bson:"content,omitempty" json:"content,omitempty"`
	Type      string             `bson:"type,omitempty" json:"type,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
