package creator_submission

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/features/creators"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const creatorSubmissionsCollection = "creatorSubmissions"

func InsertCreatorSubmission(website string, addedBy *creators.AddedByInfo) interface{} {
	newSubmission := CreatorSubmission{
		ID:        primitive.NewObjectID(),
		Website:   website,
		AddedBy:   addedBy,
		CreatedAt: time.Now(),
	}

	insertedID, err := utils.InsertDocument(creatorSubmissionsCollection, newSubmission)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllCreatorSubmissions() ([]CreatorSubmission, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(creatorSubmissionsCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var submissions []CreatorSubmission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return submissions, nil
}

func DeleteCreatorSubmissionByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(creatorSubmissionsCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error deleting creator submission for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no creator submission found with ID %s", id.Hex())
	}

	return nil
}