package module

import (
	"fmt"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertWeb3Tool(name, description, category string, chains []string, imageUrl, website, twitter, discord, telegram string) interface{} {
    newTool := models.Web3Tool{
        ID:          primitive.NewObjectID(),
        Name:        name,
        Description: description,
        Category:    category,
        Chains:      chains,
        ImageURL:    imageUrl,
        Website:     website,
        Twitter:     twitter,
        Discord:     discord,
        Telegram:    telegram,
    }

    insertedID, err := InsertDocument("web3Tool", newTool)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    return insertedID
}

func GetAllWeb3Tool() ([]models.Web3Tool, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("web3Tool")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var tools []models.Web3Tool
	if err = cursor.All(ctx, &tools); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return tools, nil
}

func GetWeb3ToolStats() (map[string]interface{}, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("web3Tool")

    pipeline := bson.A{
        bson.M{
            "$facet": bson.M{
                "total": bson.A{
                    bson.M{"$count": "count"},
                },
                "categories": bson.A{
                    bson.M{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
                },
                "chains": bson.A{
                    bson.M{"$unwind": "$chains"},
                    bson.M{"$group": bson.M{"_id": "$chains", "count": bson.M{"$sum": 1}}},
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
        "chains":     map[string]int{},
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

        chains := make(map[string]int)
        if chainArr, ok := facet["chains"].(bson.A); ok {
            for _, item := range chainArr {
                if doc, ok := item.(bson.M); ok {
                    key := ""
                    if doc["_id"] != nil {
                        key = doc["_id"].(string)
                    }
                    if count, ok := doc["count"].(int32); ok {
                        chains[key] = int(count)
                    }
                }
            }
        }
        stats["chains"] = chains
    }

    return stats, nil
}

func GetWeb3ToolByID(id primitive.ObjectID) (*models.Web3Tool, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("web3Tool")
	filter := bson.M{"_id": id}

	var result models.Web3Tool
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateWeb3ToolByID(id primitive.ObjectID, updateData models.Web3Tool) (*models.Web3Tool, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("web3Tool")

	update := bson.M{
		"$set": bson.M{
			"name":        updateData.Name,
			"description": updateData.Description,
			"category":    updateData.Category,
			"chains":      updateData.Chains,
			"imageUrl":    updateData.ImageURL,
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

func DeleteWeb3ToolByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("web3Tool")
    filter := bson.M{"_id": id}

    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("error deleting web3 tool for ID %s: %s", id.Hex(), err.Error())
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("no web3 tool found with ID %s", id.Hex())
    }

    return nil
}