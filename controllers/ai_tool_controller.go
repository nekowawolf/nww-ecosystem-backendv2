package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateAIToolCache() {
	utils.InvalidateCache("aitool", "aitool_stats")
}

func GetAllAITool(c *fiber.Ctx) error {
	tools, err := utils.GetOrSetCache("aitool", 24*time.Hour, func() ([]models.AITool, error) {
		return module.GetAllAITool()
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
	stats, err := utils.GetOrSetCache("aitool_stats", 24*time.Hour, func() (map[string]interface{}, error) {
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

func GetAIToolByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	tool, err := module.GetAIToolByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "AITool not found",
		})
	}

	return c.JSON(tool)
}

func InsertAITool(c *fiber.Ctx) error {
	var req models.AITool

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := module.InsertAITool(
		req.Name,
		req.Description,
		req.Categories,
		req.ImgURL,
		req.Website,
		req.Twitter,
		req.Discord,
		req.Telegram,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert AITool",
		})
	}

	invalidateAIToolCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "AITool created successfully",
		"insertedID": insertedID,
	})
}

func UpdateAIToolByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req models.AITool

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := models.AITool{
		Name:        req.Name,
		Description: req.Description,
		Categories:  req.Categories,
		ImgURL:      req.ImgURL,
		Website:     req.Website,
		Twitter:     req.Twitter,
		Discord:     req.Discord,
		Telegram:    req.Telegram,
	}

	updatedTool, err := module.UpdateAIToolByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "AITool not found or could not be updated",
		})
	}

	invalidateAIToolCache()
	return c.JSON(fiber.Map{
		"message": "AITool updated successfully",
		"data":    updatedTool,
	})
}

func DeleteAIToolByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = module.DeleteAIToolByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateAIToolCache()
	return c.JSON(fiber.Map{
		"message": "AITool deleted successfully",
	})
}