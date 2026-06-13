package utils

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseObjectID(c *fiber.Ctx, paramName string) (primitive.ObjectID, error) {
	idParam := c.Params(paramName)
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
		return primitive.NilObjectID, err
	}
	return id, nil
}

func ParseBody(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
		return err
	}
	return nil
}
