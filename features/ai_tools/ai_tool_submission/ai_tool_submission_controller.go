package ai_tool_submission

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/features/ai_tools"
	"github.com/nekowawolf/airdropv2/utils"
)

func SubmitAIToolHandler(c *fiber.Ctx) error {
	ip := c.IP()
	key := fmt.Sprintf("rate:ai_tool_submission:%s", ip)
	ctx := context.Background()

	count, err := config.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error during rate limiting",
		})
	}

	if count == 1 {
		config.RedisClient.Expire(ctx, key, 5*time.Minute)
	}

	if count > 10 {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Too many requests",
		})
	}

	var req AIToolSubmissionRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	if req.Website == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Website is required",
		})
	}

	if _, err := url.ParseRequestURI(req.Website); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL format for Website",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name (Added by) is required",
		})
	}

	if req.Link != "" {
		if _, err := url.ParseRequestURI(req.Link); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid URL format for Added by Link",
			})
		}
	}

	if !utils.VerifyTurnstile(req.TurnstileToken) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Invalid or missing Turnstile token",
		})
	}

	addedBy := &ai_tools.AddedByInfo{
		Name: req.Name,
		URL:  req.Link,
	}

	insertedID := InsertAIToolSubmission(req.Website, addedBy)
	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit AI Tool",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "AI Tool submitted successfully",
		"insertedID": insertedID,
	})
}

func GetAllAIToolSubmissionsHandler(c *fiber.Ctx) error {
	submissions, err := GetAllAIToolSubmissions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    submissions,
	})
}

func DeleteAIToolSubmissionHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	if err = DeleteAIToolSubmissionByID(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "AI Tool submission deleted successfully",
	})
}