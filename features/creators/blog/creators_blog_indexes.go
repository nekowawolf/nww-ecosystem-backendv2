package blog

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	createCreatorsBlogIndexes(config.Database)
}

func createCreatorsBlogIndexes(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: -1}},
		},
	}
	if _, err := db.Collection(creatorsBlogCollection).Indexes().CreateMany(ctx, indexes); err != nil {
		fmt.Printf("Failed to create indexes for creatorsBlogs: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Indexes for creatorsBlogs created successfully")
}
