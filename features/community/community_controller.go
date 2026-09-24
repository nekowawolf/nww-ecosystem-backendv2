package community

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateCommunityCache() {
	utils.InvalidateCache("cryptocommunity", "cryptocommunity_stats")
}

func GetAllCryptoCommunityHandler(c *fiber.Ctx) error {
	cryptoCommunities, err := utils.GetOrSetCache("cryptocommunity", 24*time.Hour, func() ([]CryptoCommunity, error) {
		return GetAllCryptoCommunity()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    cryptoCommunities,
	})
}

func GetCryptoCommunityStatsHandler(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("cryptocommunity_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return GetCryptoCommunityStats()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Stats retrieved successfully",
		"data":    stats,
	})
}

func GetCryptoCommunityByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	cryptoCommunity, err := GetCryptoCommunityByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "CryptoCommunity not found",
		})
	}

	return c.JSON(cryptoCommunity)
}

func InsertCryptoCommunityHandler(c *fiber.Ctx) error {
	var req CryptoCommunity

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := InsertCryptoCommunity(req)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert CryptoCommunity",
		})
	}

	invalidateCommunityCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "CryptoCommunity created successfully",
		"insertedID": insertedID,
	})
}

func UpdateCryptoCommunityByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req CryptoCommunity

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := CryptoCommunity{
		Name:        req.Name,
		Description: req.Description,
		Platforms:   req.Platforms,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Website:     req.Website,
		Link:        req.Link,
		Socials:     req.Socials,
		AddedBy:     req.AddedBy,
	}

	updatedCommunity, err := UpdateCryptoCommunityByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "CryptoCommunity not found or could not be updated",
		})
	}

	invalidateCommunityCache()
	return c.JSON(fiber.Map{
		"message": "CryptoCommunity updated successfully",
		"data":    updatedCommunity,
	})
}

func DeleteCryptoCommunityByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = DeleteCryptoCommunityByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateCommunityCache()
	return c.JSON(fiber.Map{
		"message": "CryptoCommunity deleted successfully",
	})
}