package module

import (
	"fmt"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"github.com/nekowawolf/airdropv2/utils"
)

func InsertAdmin(username, password string) (interface{}, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	newAdmin := models.Admin{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(hashedPassword),
	}

	insertedID, err := InsertDocument("admin", newAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to insert admin: %v", err)
	}

	fmt.Printf("Inserted new admin with ID: %v\n", insertedID)
	return insertedID, nil
}

func LoginAdmin(username, password string) (bool, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection("admin")

	var admin models.Admin
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&admin)
	if err == mongo.ErrNoDocuments {
		return false, fmt.Errorf("admin not found")
	} else if err != nil {
		return false, fmt.Errorf("error finding admin: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return false, fmt.Errorf("invalid password")
	}

	fmt.Println("login successful")
	return true, nil
}