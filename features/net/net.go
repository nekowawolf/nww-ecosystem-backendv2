package net

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertNet(name, description string, imageURL string, website string, categories []string, media NetMedia, socials NetSocials, addedBy *AddedByInfo) interface{} {
	newNet := Net{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
		Website:     website,
		Categories:  categories,
		Media:       media,
		Socials:     socials,
		AddedBy:     addedBy,
		CreatedAt:   time.Now(),
	}

	insertedID, err := utils.InsertDocument("net", newNet)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllNet() ([]Net, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("net")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var nets []Net
	if err = cursor.All(ctx, &nets); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return nets, nil
}

func GetNetStats() (map[string]interface{}, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("net")

	pipeline := bson.A{
		bson.M{
			"$facet": bson.M{
				"total": bson.A{
					bson.M{"$count": "count"},
				},
				"categories": bson.A{
					bson.M{"$unwind": "$categories"},
					bson.M{"$group": bson.M{"_id": "$categories", "count": bson.M{"$sum": 1}}},
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

func GetNetByID(id primitive.ObjectID) (*Net, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("net")
	filter := bson.M{"_id": id}

	var result Net
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateNetByID(id primitive.ObjectID, updateData Net) (*Net, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("net")

	setData := bson.M{
		"name":        updateData.Name,
		"description": updateData.Description,
		"image_url":   updateData.ImageURL,
		"website":     updateData.Website,
		"categories":  updateData.Categories,
		"media":       updateData.Media,
		"socials":     updateData.Socials,
	}
	if updateData.AddedBy != nil {
		setData["added_by"] = updateData.AddedBy
	}

	update := bson.M{"$set": setData}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, fmt.Errorf("error updating document: %v", err)
	}

	return &updateData, nil
}

func DeleteNetByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("net")
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting net for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no net found with ID %s", id.Hex())
	}

	return nil
}