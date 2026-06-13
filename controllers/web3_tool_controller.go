package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateWeb3ToolCache() {
	utils.InvalidateCache("web3tool", "web3tool_stats")
}

func GetAllWeb3Tool(c *fiber.Ctx) error {
	tools, err := utils.GetOrSetCache("web3tool", 24*time.Hour, func() ([]models.Web3Tool, error) {
		return module.GetAllWeb3Tool()
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

func GetWeb3ToolStats(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("web3tool_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return module.GetWeb3ToolStats()
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

func GetWeb3ToolByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	tool, err := module.GetWeb3ToolByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Web3Tool not found",
		})
	}

	return c.JSON(tool)
}

func InsertWeb3Tool(c *fiber.Ctx) error {
	var req models.Web3Tool

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := module.InsertWeb3Tool(
		req.Name,
		req.Description,
		req.Category,
		req.Chains,
		req.ImageURL,
		req.Website,
		req.Twitter,
		req.Discord,
		req.Telegram,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert Web3Tool",
		})
	}

	invalidateWeb3ToolCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Web3Tool created successfully",
		"insertedID": insertedID,
	})
}

func UpdateWeb3ToolByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req models.Web3Tool

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := models.Web3Tool{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Chains:      req.Chains,
		ImageURL:    req.ImageURL,
		Website:     req.Website,
		Twitter:     req.Twitter,
		Discord:     req.Discord,
		Telegram:    req.Telegram,
	}

	updatedTool, err := module.UpdateWeb3ToolByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Web3Tool not found or could not be updated",
		})
	}

	invalidateWeb3ToolCache()
	return c.JSON(fiber.Map{
		"message": "Web3Tool updated successfully",
		"data":    updatedTool,
	})
}

func DeleteWeb3ToolByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = module.DeleteWeb3ToolByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateWeb3ToolCache()
	return c.JSON(fiber.Map{
		"message": "Web3Tool deleted successfully",
	})
}