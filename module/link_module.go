package module

import (
	"github.com/nekowawolf/airdropv2/utils"
	"errors"
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ==================== PROFILE CRUD ====================

func GetProfile() (*models.Profile, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("profile")
	var profile models.Profile

	err := collection.FindOne(ctx, bson.M{}).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func UpdateProfile(profile models.Profile) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("profile")

	var existing models.Profile
	err := collection.FindOne(ctx, bson.M{}).Decode(&existing)

	if err != nil {
		profile.ID = primitive.NewObjectID()
		_, err = collection.InsertOne(ctx, profile)
		return err
	}

	profile.ID = existing.ID
	filter := bson.M{"_id": existing.ID}
	update := bson.M{"$set": profile}
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// ==================== POSTS CRUD ====================

func GetAllPosts() ([]models.LinkPost, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("link_posts")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("GetAllPosts Find: %v", err)
	}
	defer cursor.Close(ctx)

	var posts []models.LinkPost
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("GetAllPosts All: %v", err)
	}

	return posts, nil
}

func GetPostsPaginated(page, limit int, category, search string) ([]models.LinkPost, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("link_posts")

	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{}

	if category != "" && category != "all" {
		filter["category"] = category
	}

	if search != "" {
		filter["caption"] = bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("GetPostsPaginated Find: %v", err)
	}
	defer cursor.Close(ctx)

	var posts []models.LinkPost
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("GetPostsPaginated All: %v", err)
	}

	// if posts is null, return empty array instead of nil
	if posts == nil {
		posts = []models.LinkPost{}
	}

	return posts, nil
}

func GetPostStats() (map[string]int64, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("link_posts")
	
	stats := make(map[string]int64)
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats["all"] = total

	aiReq, _ := collection.CountDocuments(ctx, bson.M{"category": "AI Prompts"})
	stats["AI Prompts"] = aiReq

	tplReq, _ := collection.CountDocuments(ctx, bson.M{"category": "Templates"})
	stats["Templates"] = tplReq

	projReq, _ := collection.CountDocuments(ctx, bson.M{"category": "projects"})
	stats["projects"] = projReq

	return stats, nil
}

func GetPostByID(id primitive.ObjectID) (*models.LinkPost, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("link_posts")
	var post models.LinkPost

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func InsertPost(post models.LinkPost) (interface{}, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()
	collection := config.Database.Collection("link_posts")

	now := time.Now()
	post.ID = primitive.NewObjectID()
	post.CreatedAt = now
	post.Views = 0
	post.IsVerified = true

	profile, err := GetProfile()
	if err == nil {
		post.Name = profile.Name
		post.Username = profile.Username
	}

	result, err := collection.InsertOne(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("InsertPost: %v", err)
	}

	return result.InsertedID, nil
}

func UpdatePost(id primitive.ObjectID, post models.LinkPost) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("link_posts")

	existing, err := GetPostByID(id)
	if err != nil {
		return errors.New("post not found")
	}

	post.ID = existing.ID
	post.CreatedAt = existing.CreatedAt
	post.Views = existing.Views
	post.IsVerified = true
	post.Name = existing.Name
	post.Username = existing.Username

	filter := bson.M{"_id": id}
	update := bson.M{"$set": post}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("UpdatePost: %v", err)
	}

	if result.ModifiedCount == 0 {
		return errors.New("no data has been changed")
	}

	return nil
}

func DeletePost(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("link_posts")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("DeletePost: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("post not found")
	}

	viewCollection := config.Database.Collection("view_stats")
	_, _ = viewCollection.DeleteMany(ctx, bson.M{"post_id": id})

	return nil
}

// ==================== VIEWS SYSTEM ====================

func IncrementPostView(postID primitive.ObjectID, sessionID string) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	viewCollection := config.Database.Collection("view_stats")

	var existingView models.ViewStats
	err := viewCollection.FindOne(ctx, bson.M{
		"post_id":    postID,
		"session_id": sessionID,
	}).Decode(&existingView)

	if err == nil {
		return nil
	}

	view := models.ViewStats{
		ID:        primitive.NewObjectID(),
		PostID:    postID,
		SessionID: sessionID,
		ViewedAt:  time.Now(),
	}

	_, err = viewCollection.InsertOne(ctx, view)
	if err != nil {
		return err
	}

	postCollection := config.Database.Collection("link_posts")
	_, err = postCollection.UpdateOne(
		ctx,
		bson.M{"_id": postID},
		bson.M{"$inc": bson.M{"views": 1}},
	)

	return err
}

func GetPostViewCount(postID primitive.ObjectID) (int64, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("view_stats")
	return collection.CountDocuments(ctx, bson.M{"post_id": postID})
}