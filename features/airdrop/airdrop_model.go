package airdrop

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddedByInfo struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	URL  string `bson:"url,omitempty" json:"url,omitempty"`
}

type AirdropPublic struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Task        string             `json:"task,omitempty" bson:"task,omitempty"`
	Website     string             `json:"website,omitempty" bson:"website,omitempty"`
	Level       string             `json:"level,omitempty" bson:"level,omitempty"`
	Status      string             `json:"status,omitempty" bson:"status,omitempty"`
	Backed      string             `json:"backed,omitempty" bson:"backed,omitempty"`
	Funds       string             `json:"funds,omitempty" bson:"funds,omitempty"`
	IsVesting   bool               `json:"is_vesting" bson:"is_vesting"`
	IsPaid      bool               `json:"is_paid" bson:"is_paid"`
	ClaimURL    string             `json:"claim_url,omitempty" bson:"claim_url,omitempty"`
	Discord     string             `json:"discord,omitempty" bson:"discord,omitempty"`
	Twitter     string             `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Telegram    string             `json:"telegram,omitempty" bson:"telegram,omitempty"`
	ImageURL    string             `json:"image_url,omitempty" bson:"image_url,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	GuideURL    string             `json:"guide_url,omitempty" bson:"guide_url,omitempty"`
	AddedBy     *AddedByInfo       `json:"added_by,omitempty" bson:"added_by,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	EndedAt     *time.Time         `json:"ended_at,omitempty" bson:"ended_at,omitempty"`
}

type AirdropAdmin struct {
	AirdropPublic `bson:",inline"`
	Supply        string  `json:"supply,omitempty" bson:"supply,omitempty"`
	Fdv           string  `json:"fdv,omitempty" bson:"fdv,omitempty"`
	MarketCap     string  `json:"market_cap,omitempty" bson:"market_cap,omitempty"`
	Price         float64 `json:"price,omitempty" bson:"price,omitempty"`
	USDIncome     int     `json:"usd_income,omitempty" bson:"usd_income,omitempty"`
}