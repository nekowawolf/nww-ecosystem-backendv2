package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMessage(c *fiber.Ctx) error {
	var message models.Message
	err := config.Database.Collection("messages").FindOne(context.Background(), bson.M{}).Decode(&message)
	
	if err != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    models.Message{},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    message,
	})
}

func UpdateMessage(c *fiber.Ctx) error {
	var body models.Message
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	body.ID = [12]byte{}

	update := bson.M{
		"$set": body,
	}
	opts := options.Update().SetUpsert(true)

	_, err := config.Database.Collection("messages").UpdateOne(context.Background(), bson.M{}, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update messages",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Messages updated successfully",
	})
}