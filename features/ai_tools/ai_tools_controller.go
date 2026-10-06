package ai_tools

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateAIToolsCache() {
	utils.InvalidateCache("aitools", "aitools_stats")
}

func GetAllAIToolsHandler(c *fiber.Ctx) error {
	tools, err := utils.GetOrSetCache("aitools", 24*time.Hour, func() ([]AITools, error) {
		return GetAllAITools()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    tools,
	})
}

func GetAIToolStatsHandler(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("aitools_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return GetAIToolStats()
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

func GetAIToolsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	tool, err := GetAIToolsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "AITools not found",
		})
	}

	return c.JSON(tool)
}

func InsertAIToolsHandler(c *fiber.Ctx) error {
	var req AITools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := InsertAITools(
		req.Name,
		req.Description,
		req.ImageURL,
		req.Website,
		req.Categories,
		req.Media,
		req.Socials,
		req.AddedBy,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert AITools",
		})
	}

	invalidateAIToolsCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "AITools created successfully",
		"insertedID": insertedID,
	})
}

func UpdateAIToolsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req AITools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := AITools{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Website:     req.Website,
		Categories:  req.Categories,
		Media:       req.Media,
		Socials:     req.Socials,
		AddedBy:     req.AddedBy,
	}

	updatedTool, err := UpdateAIToolsByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "AITools not found or could not be updated",
		})
	}

	invalidateAIToolsCache()
	return c.JSON(fiber.Map{
		"message": "AITools updated successfully",
		"data":    updatedTool,
	})
}

func DeleteAIToolsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = DeleteAIToolsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateAIToolsCache()
	return c.JSON(fiber.Map{
		"message": "AITools deleted successfully",
	})
}