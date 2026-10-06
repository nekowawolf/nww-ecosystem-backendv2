package guild_submission

import (
	"time"

	"github.com/nekowawolf/airdropv2/features/guild"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GuildSubmissionRequest struct {
	GuildLink      string `json:"guild_link"`
	Name           string `json:"name"`
	Link           string `json:"link,omitempty"`
	TurnstileToken string `json:"turnstile_token"`
}

type GuildSubmission struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	GuildLink string             `bson:"guild_link,omitempty" json:"guild_link,omitempty"`
	AddedBy   *guild.AddedByInfo `bson:"added_by,omitempty" json:"added_by,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}