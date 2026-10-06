package creators

import (
	"fmt"
	"time"
	"github.com/nekowawolf/airdropv2/config"

	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertCreators(name, description string, imageURL string, website string, category, language string, openToWork bool, socials CreatorsSocials, platforms CreatorsPlatforms, addedBy *AddedByInfo) interface{} {
    newCreator := Creators{
        ID:          primitive.NewObjectID(),
        Name:        name,
        Description: description,
        ImageURL:    imageURL,
        Website:     website,
        Category:    category,
        Language:    language,
        OpenToWork:  openToWork,
        Socials:     socials,
        Platforms:   platforms,
        AddedBy:     addedBy,
        CreatedAt:   time.Now(),
    }

    insertedID, err := utils.InsertDocument("creators", newCreator)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    return insertedID
}

func GetAllCreators() ([]Creators, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("creators")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var creators []Creators
	if err = cursor.All(ctx, &creators); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return creators, nil
}

func GetCreatorsStats() (map[string]interface{}, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("creators")

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
                    if key != "" {
                        if count, ok := doc["count"].(int32); ok {
                            categories[key] = int(count)
                        }
                    }
                }
            }
        }
        stats["categories"] = categories
    }

    return stats, nil
}

func GetCreatorsByID(id primitive.ObjectID) (*Creators, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("creators")
	filter := bson.M{"_id": id}

	var result Creators
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateCreatorsByID(id primitive.ObjectID, updateData Creators) (*Creators, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("creators")

	setData := bson.M{
		"name":          updateData.Name,
		"description":   updateData.Description,
		"image_url":     updateData.ImageURL,
		"website":       updateData.Website,
		"category":      updateData.Category,
		"language":      updateData.Language,
		"open_to_work":  updateData.OpenToWork,
		"socials":       updateData.Socials,
		"platforms":     updateData.Platforms,
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

func DeleteCreatorsByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

    collection := config.Database.Collection("creators")
    filter := bson.M{"_id": id}

    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return fmt.Errorf("error deleting creator for ID %s: %s", id.Hex(), err.Error())
    }

    if result.DeletedCount == 0 {
        return fmt.Errorf("no creator found with ID %s", id.Hex())
    }

    return nil
}