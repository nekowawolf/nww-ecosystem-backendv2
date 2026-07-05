package module

import (
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertNote(title, content, noteType string) interface{} {
	newNote := models.Notes{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		Type:      noteType,
		CreatedAt: time.Now(),
	}

	insertedID, err := InsertDocument("notes", newNote)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return insertedID
}

func GetAllNotes() ([]models.Notes, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("notes")
	
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
    
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %v", err)
	}
	defer cursor.Close(ctx)

	var notes []models.Notes
	if err = cursor.All(ctx, &notes); err != nil {
		return nil, fmt.Errorf("error decoding data: %v", err)
	}

	return notes, nil
}

func GetNoteByID(id primitive.ObjectID) (*models.Notes, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("notes")
	filter := bson.M{"_id": id}

	var result models.Notes
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateNoteByID(id primitive.ObjectID, updateData models.Notes) (*models.Notes, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("notes")

	update := bson.M{
		"$set": bson.M{
			"title":   updateData.Title,
			"content": updateData.Content,
			"type":    updateData.Type,
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, fmt.Errorf("error updating document: %v", err)
	}

	return &updateData, nil
}

func DeleteNoteByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("notes")
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting note for ID %s: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no note found with ID %s", id.Hex())
	}

	return nil
}
