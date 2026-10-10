package airdrop_submission

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/features/airdrop"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const airdropSubmissionsCollection = "airdropSubmissions"

func InsertAirdropSubmission(website string, addedBy *airdrop.AddedByInfo) interface{} {
	newSubmission := AirdropSubmission{
		ID:        primitive.NewObjectID(),
		Website:   website,
		AddedBy:   addedBy,
		CreatedAt: time.Now(),
	}

	insertedID, err := utils.InsertDocument(airdropSubmissionsCollection, newSubmission)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllAirdropSubmissions() ([]AirdropSubmission, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(airdropSubmissionsCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var submissions []AirdropSubmission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return submissions, nil
}

func DeleteAirdropSubmissionByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(airdropSubmissionsCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error deleting airdrop submission for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no airdrop submission found with ID %s", id.Hex())
	}

	return nil
}