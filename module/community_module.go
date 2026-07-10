package module

import (
	"fmt"
	"time"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertCryptoCommunity(name, platforms, category, imgURL, linkURL string) interface{} {
    newCrypto := models.CryptoCommunity{
        ID:        primitive.NewObjectID(),
        Name:      name,
        Platforms: platforms,
        Category:  category,
        ImgURL:    imgURL,
        LinkURL:   linkURL,
        CreatedAt: time.Now(),
    }

    insertedID, err := InsertDocument("cryptoCommunity", newCrypto)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    fmt.Printf("Inserted new crypto community with ID: %v\n", insertedID)
    return insertedID
}

func GetAllCryptoCommunity() ([]models.CryptoCommunity, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("cryptoCommunity")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var communities []models.CryptoCommunity
	if err = cursor.All(ctx, &communities); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return communities, nil
}

func GetCryptoCommunityStats() (map[string]interface{}, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("cryptoCommunity")

    pipeline := bson.A{
        bson.M{
            "$facet": bson.M{
                "total": bson.A{
                    bson.M{"$count": "count"},
                },
                "categories": bson.A{
                    bson.M{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
                },
                "platforms": bson.A{
                    bson.M{"$group": bson.M{"_id": "$platforms", "count": bson.M{"$sum": 1}}},
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
        "platforms":  map[string]int{},
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

        platforms := make(map[string]int)
        if platArr, ok := facet["platforms"].(bson.A); ok {
            for _, item := range platArr {
                if doc, ok := item.(bson.M); ok {
                    key := ""
                    if doc["_id"] != nil {
                        key = doc["_id"].(string)
                    }
                    if count, ok := doc["count"].(int32); ok {
                        platforms[key] = int(count)
                    }
                }
            }
        }
        stats["platforms"] = platforms
    }

    return stats, nil
}

func GetCryptoCommunityByID(id primitive.ObjectID) (*models.CryptoCommunity, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("cryptoCommunity")
	filter := bson.M{"_id": id}

	var result models.CryptoCommunity
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateCryptoCommunityByID(id primitive.ObjectID, updateData models.CryptoCommunity) (*models.CryptoCommunity, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("cryptoCommunity")

	update := bson.M{
		"$set": bson.M{
			"name":       updateData.Name,
			"platforms":  updateData.Platforms,
			"category":   updateData.Category,
			"img_url":    updateData.ImgURL,
			"link_url":   updateData.LinkURL,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, fmt.Errorf("error updating document: %v", err)
	}

	return &updateData, nil
}

func DeleteCryptoCommunityByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("cryptoCommunity")
    filter := bson.M{"_id": id}

    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("error deleting crypto community for ID %s: %s", id.Hex(), err.Error())
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("no crypto community found with ID %s", id.Hex())
    }

    return nil
}