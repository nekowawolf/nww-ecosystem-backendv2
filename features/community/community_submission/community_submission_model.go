package community_submission

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddedByInfo struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	URL  string `bson:"url,omitempty" json:"url,omitempty"`
}

type CommunitySubmissionRequest struct {
	CommunityLink  string `json:"community_link"`
	Name           string `json:"name"`
	Link           string `json:"link,omitempty"`
	TurnstileToken string `json:"turnstile_token"`
}

type CommunitySubmission struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CommunityLink string             `bson:"community_link,omitempty" json:"community_link,omitempty"`
	AddedBy       *AddedByInfo       `bson:"added_by,omitempty" json:"added_by,omitempty"`
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}