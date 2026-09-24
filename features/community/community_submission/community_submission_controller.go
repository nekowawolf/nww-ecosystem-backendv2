package community_submission

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/utils"
)

func SubmitCommunityHandler(c *fiber.Ctx) error {
	ip := c.IP()
	key := fmt.Sprintf("rate:community_submission:%s", ip)
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

	var req CommunitySubmissionRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	if req.CommunityLink == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Community Link is required",
		})
	}

	if _, err := url.ParseRequestURI(req.CommunityLink); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL format for Community Link",
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

	addedBy := &AddedByInfo{
		Name: req.Name,
		URL:  req.Link,
	}

	insertedID := InsertCommunitySubmission(&req, addedBy)
	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit community",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Community submitted successfully",
		"insertedID": insertedID,
	})
}

func GetAllCommunitySubmissionsHandler(c *fiber.Ctx) error {
	submissions, err := GetAllCommunitySubmissions()
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

func DeleteCommunitySubmissionHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = DeleteCommunitySubmissionByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Community submission deleted successfully",
	})
}