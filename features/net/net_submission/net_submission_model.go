package net_submission

import (
	"time"

	"github.com/nekowawolf/airdropv2/features/net"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NetSubmissionRequest struct {
	Website        string `json:"website"`
	Name           string `json:"name"`
	Link           string `json:"link,omitempty"`
	TurnstileToken string `json:"turnstile_token"`
}

type NetSubmission struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Website   string             `bson:"website,omitempty" json:"website,omitempty"`
	AddedBy   *net.AddedByInfo   `bson:"added_by,omitempty" json:"added_by,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}