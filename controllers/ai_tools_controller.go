package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateAIToolsCache() {
	utils.InvalidateCache("aitools", "aitools_stats")
}

func GetAllAITools(c *fiber.Ctx) error {
	tools, err := utils.GetOrSetCache("aitools", 24*time.Hour, func() ([]models.AITools, error) {
		return module.GetAllAITools()
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

func GetAIToolStats(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("aitools_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return module.GetAIToolStats()
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

func GetAIToolsByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	tool, err := module.GetAIToolsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "AITools not found",
		})
	}

	return c.JSON(tool)
}

func InsertAITools(c *fiber.Ctx) error {
	var req models.AITools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := module.InsertAITools(
		req.Name,
		req.Description,
		req.Categories,
		req.VideoURL,
		req.ImgURL,
		req.Website,
		req.Twitter,
		req.Instagram,
		req.Discord,
		req.Telegram,
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

func UpdateAIToolsByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req models.AITools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := models.AITools{
		Name:        req.Name,
		Description: req.Description,
		Categories:  req.Categories,
		VideoURL:    req.VideoURL,
		ImgURL:      req.ImgURL,
		Website:     req.Website,
		Twitter:     req.Twitter,
		Instagram:   req.Instagram,
		Discord:     req.Discord,
		Telegram:    req.Telegram,
	}

	updatedTool, err := module.UpdateAIToolsByID(id, updateData)
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

func DeleteAIToolsByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = module.DeleteAIToolsByID(id)
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