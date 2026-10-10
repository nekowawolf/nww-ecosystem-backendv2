package airdrop_submission

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/features/airdrop"
	"github.com/nekowawolf/airdropv2/utils"
)

func SubmitAirdropHandler(c *fiber.Ctx) error {
	ip := c.IP()
	key := fmt.Sprintf("rate:airdrop_submission:%s", ip)
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

	var req AirdropSubmissionRequest
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

	addedBy := &airdrop.AddedByInfo{
		Name: req.Name,
		URL:  req.Link,
	}

	insertedID := InsertAirdropSubmission(req.Website, addedBy)
	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit Airdrop",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Airdrop submitted successfully",
		"insertedID": insertedID,
	})
}

func GetAllAirdropSubmissionsHandler(c *fiber.Ctx) error {
	submissions, err := GetAllAirdropSubmissions()
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

func DeleteAirdropSubmissionHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	if err = DeleteAirdropSubmissionByID(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Airdrop submission deleted successfully",
	})
}