package module

import (
	"fmt"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertAITool(name, description string, categories []string, imgURL, website, twitter, discord, telegram string) interface{} {
    newTool := models.AITool{
        ID:          primitive.NewObjectID(),
        Name:        name,
        Description: description,
        Categories:  categories,
        ImgURL:      imgURL,
        Website:     website,
        Twitter:     twitter,
        Discord:     discord,
        Telegram:    telegram,
    }

    insertedID, err := InsertDocument("aiTool", newTool)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    return insertedID
}

func GetAllAITool() ([]models.AITool, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("aiTool")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var tools []models.AITool
	if err = cursor.All(ctx, &tools); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return tools, nil
}

func GetAIToolStats() (map[string]interface{}, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("aiTool")

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

func GetAIToolByID(id primitive.ObjectID) (*models.AITool, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("aiTool")
	filter := bson.M{"_id": id}

	var result models.AITool
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateAIToolByID(id primitive.ObjectID, updateData models.AITool) (*models.AITool, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("aiTool")

	update := bson.M{
		"$set": bson.M{
			"name":        updateData.Name,
			"description": updateData.Description,
			"categories":  updateData.Categories,
			"imgURL":      updateData.ImgURL,
			"website":     updateData.Website,
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

func DeleteAIToolByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("aiTool")
    filter := bson.M{"_id": id}

    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("error deleting ai tool for ID %s: %s", id.Hex(), err.Error())
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("no ai tool found with ID %s", id.Hex())
    }

    return nil
}