package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GithubRepo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	RepoURL     string             `bson:"repo_url,omitempty" json:"repo_url,omitempty"`
	Owner       string             `bson:"owner,omitempty" json:"owner,omitempty"`
	RepoName    string             `bson:"repo_name,omitempty" json:"repo_name,omitempty"`
	Twitter     string             `bson:"twitter,omitempty" json:"twitter,omitempty"`
	Discord     string             `bson:"discord,omitempty" json:"discord,omitempty"`
	Telegram    string             `bson:"telegram,omitempty" json:"telegram,omitempty"`
}