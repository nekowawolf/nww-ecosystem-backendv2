package community_submission

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"

	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertCommunitySubmission(req *CommunitySubmissionRequest, addedBy *AddedByInfo) interface{} {
	newSubmission := CommunitySubmission{
		ID:            primitive.NewObjectID(),
		CommunityLink: req.CommunityLink,
		AddedBy:       addedBy,
		CreatedAt:     time.Now(),
	}

	insertedID, err := utils.InsertDocument("communitySubmissions", newSubmission)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllCommunitySubmissions() ([]CommunitySubmission, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("communitySubmissions")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var submissions []CommunitySubmission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return submissions, nil
}

func DeleteCommunitySubmissionByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("communitySubmissions")
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting community submission for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no community submission found with ID %s", id.Hex())
	}

	return nil
}