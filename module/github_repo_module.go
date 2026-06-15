package module

import (
	"fmt"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertGithubRepo(name, description, category, repoURL, owner, repoName, twitter, discord, telegram string) interface{} {
	newRepo := models.GithubRepo{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		Category:    category,
		RepoURL:     repoURL,
		Owner:       owner,
		RepoName:    repoName,
		Twitter:     twitter,
		Discord:     discord,
		Telegram:    telegram,
	}

	insertedID, err := InsertDocument("githubRepos", newRepo)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllGithubRepos() ([]models.GithubRepo, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("githubRepos")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var repos []models.GithubRepo
	if err = cursor.All(ctx, &repos); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return repos, nil
}

func GetGithubRepoStats() (map[string]interface{}, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("githubRepos")

	pipeline := bson.A{
		bson.M{
			"$facet": bson.M{
				"total": bson.A{
					bson.M{"$count": "count"},
				},
				"categories": bson.A{
					bson.M{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
				},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error aggregating data: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("error decoding aggregation: %v", err)
	}

	stats := map[string]interface{}{
		"total":      0,
		"categories": map[string]int{},
	}

	if len(results) > 0 {
		facet := results[0]

		if totalArr, ok := facet["total"].(bson.A); ok && len(totalArr) > 0 {
			if totalDoc, ok := totalArr[0].(bson.M); ok {
				if count, ok := totalDoc["count"].(int32); ok {
					stats["total"] = int(count)
				}
			}
		}

		categories := make(map[string]int)
		if catArr, ok := facet["categories"].(bson.A); ok {
			for _, item := range catArr {
				if doc, ok := item.(bson.M); ok {
					key := ""
					if doc["_id"] != nil {
						key = doc["_id"].(string)
					}
					if count, ok := doc["count"].(int32); ok {
						categories[key] = int(count)
					}
				}
			}
		}
		stats["categories"] = categories
	}

	return stats, nil
}

func GetGithubRepoByID(id primitive.ObjectID) (*models.GithubRepo, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("githubRepos")
	filter := bson.M{"_id": id}

	var result models.GithubRepo
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateGithubRepoByID(id primitive.ObjectID, updateData models.GithubRepo) (*models.GithubRepo, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("githubRepos")

	update := bson.M{
		"$set": bson.M{
			"name":        updateData.Name,
			"description": updateData.Description,
			"category":    updateData.Category,
			"repo_url":    updateData.RepoURL,
			"owner":       updateData.Owner,
			"repo_name":   updateData.RepoName,
			"twitter":     updateData.Twitter,
			"discord":     updateData.Discord,
			"telegram":    updateData.Telegram,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, fmt.Errorf("error updating document: %v", err)
	}

	return &updateData, nil
}

func DeleteGithubRepoByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("githubRepos")
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting github repo for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no github repo found with ID %s", id.Hex())
	}

	return nil
}