package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func invalidateCommunityCache() {
	utils.InvalidateCache("cryptocommunity", "cryptocommunity_stats")
}

func GetAllCryptoCommunity(c *fiber.Ctx) error {
	cryptoCommunities, err := utils.GetOrSetCache("cryptocommunity", 24*time.Hour, func() ([]models.CryptoCommunity, error) {
		return module.GetAllCryptoCommunity()
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

func GetCryptoCommunityStats(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("cryptocommunity_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return module.GetCryptoCommunityStats()
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

func GetCryptoCommunityByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	cryptoCommunity, err := module.GetCryptoCommunityByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "CryptoCommunity not found",
		})
	}

	return c.JSON(cryptoCommunity)
}

func InsertCryptoCommunity(c *fiber.Ctx) error {
	var req models.CryptoCommunity

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	insertedID := module.InsertCryptoCommunity(
		req.Name,
		req.Platforms,
		req.Category,
		req.ImgURL,
		req.LinkURL,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert CryptoCommunity",
		})
	}

	invalidateCommunityCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "CryptoCommunity created successfully",
		"insertedID": insertedID,
	})
}

func GetCryptoCommunityByName(c *fiber.Ctx) error {
    name := c.Params("name")
    
    data, err := module.GetCryptoCommunityByName(name)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to retrieve CryptoCommunity by Name",
        })
    }
    
    if len(data) == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "No CryptoCommunity found with the specified name",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Data retrieved successfully",
        "data":    data,
    })
}

func UpdateCryptoCommunityByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var req models.CryptoCommunity

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updateData := models.CryptoCommunity{
		Name:      req.Name,
		Platforms: req.Platforms,
		Category:  req.Category,
		ImgURL:    req.ImgURL,
		LinkURL:   req.LinkURL,
	}

	updatedCommunity, err := module.UpdateCryptoCommunityByID(id, updateData)
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

func DeleteCryptoCommunityByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	err = module.DeleteCryptoCommunityByID(id)
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