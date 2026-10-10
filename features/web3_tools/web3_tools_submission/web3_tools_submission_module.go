package web3_tools_submission

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/features/web3_tools"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const web3ToolsSubmissionsCollection = "web3ToolsSubmissions"

func InsertWeb3ToolsSubmission(website string, addedBy *web3_tools.AddedByInfo) interface{} {
	newSubmission := Web3ToolsSubmission{
		ID:        primitive.NewObjectID(),
		Website:   website,
		AddedBy:   addedBy,
		CreatedAt: time.Now(),
	}

	insertedID, err := utils.InsertDocument(web3ToolsSubmissionsCollection, newSubmission)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllWeb3ToolsSubmissions() ([]Web3ToolsSubmission, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(web3ToolsSubmissionsCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var submissions []Web3ToolsSubmission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return submissions, nil
}

func DeleteWeb3ToolsSubmissionByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(web3ToolsSubmissionsCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error deleting Web3 Tools submission for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no Web3 Tools submission found with ID %s", id.Hex())
	}

	return nil
}