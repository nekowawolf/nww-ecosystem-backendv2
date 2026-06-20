package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateWeb3ToolsCache() {
	utils.InvalidateCache("web3tools", "web3tools_stats")
}

func GetAllWeb3Tools(c *fiber.Ctx) error {
	tools, err := utils.GetOrSetCache("web3tools", 24*time.Hour, func() ([]models.Web3Tools, error) {
		return module.GetAllWeb3Tools()
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
	stats, err := utils.GetOrSetCache("web3tools_stats", 24*time.Hour, func() (map[string]interface{}, error) {
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

func GetWeb3ToolsByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	tool, err := module.GetWeb3ToolsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Web3Tools not found",
		})
	}

	return c.JSON(tool)
}

func InsertWeb3Tools(c *fiber.Ctx) error {
	var req models.Web3Tools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := module.InsertWeb3Tools(
		req.Name,
		req.Description,
		req.Category,
		req.Chains,
		req.ImageURL,
		req.Website,
		req.Twitter,
		req.Instagram,
		req.Discord,
		req.Telegram,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert Web3Tools",
		})
	}

	invalidateWeb3ToolsCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Web3Tools created successfully",
		"insertedID": insertedID,
	})
}

func UpdateWeb3ToolsByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req models.Web3Tools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := models.Web3Tools{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Chains:      req.Chains,
		ImageURL:    req.ImageURL,
		Website:     req.Website,
		Twitter:     req.Twitter,
		Instagram:   req.Instagram,
		Discord:     req.Discord,
		Telegram:    req.Telegram,
	}

	updatedTool, err := module.UpdateWeb3ToolsByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Web3Tools not found or could not be updated",
		})
	}

	invalidateWeb3ToolsCache()
	return c.JSON(fiber.Map{
		"message": "Web3Tools updated successfully",
		"data":    updatedTool,
	})
}

func DeleteWeb3ToolsByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = module.DeleteWeb3ToolsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateWeb3ToolsCache()
	return c.JSON(fiber.Map{
		"message": "Web3Tools deleted successfully",
	})
}