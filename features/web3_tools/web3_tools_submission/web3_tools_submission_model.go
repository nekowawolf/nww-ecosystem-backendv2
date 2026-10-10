package web3_tools_submission

import (
	"time"

	"github.com/nekowawolf/airdropv2/features/web3_tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Web3ToolsSubmissionRequest struct {
	Website        string `json:"website"`
	Name           string `json:"name"`
	Link           string `json:"link,omitempty"`
	TurnstileToken string `json:"turnstile_token"`
}

type Web3ToolsSubmission struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty" json:"_id,omitempty"`
	Website   string                  `bson:"website,omitempty" json:"website,omitempty"`
	AddedBy   *web3_tools.AddedByInfo `bson:"added_by,omitempty" json:"added_by,omitempty"`
	CreatedAt time.Time               `bson:"created_at,omitempty" json:"created_at,omitempty"`
}