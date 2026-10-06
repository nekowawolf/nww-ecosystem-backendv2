package ai_tool_submission

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/features/ai_tools"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const aiToolSubmissionsCollection = "aiToolSubmissions"

func InsertAIToolSubmission(website string, addedBy *ai_tools.AddedByInfo) interface{} {
	newSubmission := AIToolSubmission{
		ID:        primitive.NewObjectID(),
		Website:   website,
		AddedBy:   addedBy,
		CreatedAt: time.Now(),
	}

	insertedID, err := utils.InsertDocument(aiToolSubmissionsCollection, newSubmission)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllAIToolSubmissions() ([]AIToolSubmission, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(aiToolSubmissionsCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var submissions []AIToolSubmission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return submissions, nil
}

func DeleteAIToolSubmissionByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(aiToolSubmissionsCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error deleting AI Tool submission for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no AI Tool submission found with ID %s", id.Hex())
	}

	return nil
}