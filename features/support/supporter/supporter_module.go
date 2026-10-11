package supporter

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"

	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllSupporters() ([]Supporter, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("supporters")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var supporters []Supporter
	if err = cursor.All(ctx, &supporters); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return supporters, nil
}

func GetSupporterByID(id primitive.ObjectID) (*Supporter, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("supporters")
	var supporter Supporter
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&supporter)
	if err != nil {
		return nil, fmt.Errorf("error retrieving supporter by ID %s: %v", id.Hex(), err)
	}

	return &supporter, nil
}

func InsertSupporter(supporter *Supporter) interface{} {
	newSupporter := Supporter{
		ID:        primitive.NewObjectID(),
		Name:      supporter.Name,
		URL:       supporter.URL,
		Platform:  supporter.Platform,
		Proof:     supporter.Proof,
		Amount:    supporter.Amount,
		CreatedAt: time.Now(),
	}

	insertedID, err := utils.InsertDocument("supporters", newSupporter)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func UpdateSupporterByID(id primitive.ObjectID, req *Supporter) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("supporters")
	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"name":     req.Name,
			"url":      req.URL,
			"platform": req.Platform,
			"proof":    req.Proof,
			"amount":   req.Amount,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating supporter with ID %s: %s", id.Hex(), err.Error())
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no supporter found with ID %s", id.Hex())
	}

	return nil
}

func DeleteSupporterByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("supporters")
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting supporter with ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no supporter found with ID %s", id.Hex())
	}

	return nil
}