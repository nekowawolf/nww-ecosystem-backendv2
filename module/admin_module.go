package module

import (
	"context"
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InsertOneDocAdmin(collection string, doc interface{}) (interface{}, error) {
	collectionRef := config.Database.Collection(collection)
	insertResult, err := collectionRef.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, fmt.Errorf("InsertOneDocAdmin: %v", err)
	}
	return insertResult.InsertedID, nil
}

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

	insertedID, err := InsertOneDocAdmin("admin", newAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to insert admin: %v", err)
	}

	fmt.Printf("Inserted new admin with ID: %v\n", insertedID)
	return insertedID, nil
}

func LoginAdmin(username, password string) (bool, error) {
	collection := config.Database.Collection("admin")

	var admin models.Admin
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&admin)
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

func SaveRefreshToken(adminID, token string, expiresAt time.Time) error {
	collection := config.Database.Collection("refresh_tokens")
	doc := models.RefreshToken{
		ID:        primitive.NewObjectID(),
		Token:     token,
		AdminID:   adminID,
		ExpiresAt: expiresAt,
	}
	_, err := collection.InsertOne(context.TODO(), doc)
	return err
}

func CheckRefreshToken(token string) bool {
	collection := config.Database.Collection("refresh_tokens")
	var rt models.RefreshToken

	err := collection.FindOne(context.TODO(), bson.M{"token": token}).Decode(&rt)
	if err != nil {
		return false
	}

	if time.Now().After(rt.ExpiresAt) {
		return false
	}

	return true
}

func DeleteRefreshToken(token string) error {
	collection := config.Database.Collection("refresh_tokens")
	_, err := collection.DeleteOne(context.TODO(), bson.M{"token": token})
	return err
}