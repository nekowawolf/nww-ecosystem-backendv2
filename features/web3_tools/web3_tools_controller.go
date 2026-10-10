package web3_tools

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateWeb3ToolsCache() {
	utils.InvalidateCache("web3tools", "web3tools_stats")
}

func GetAllWeb3ToolsHandler(c *fiber.Ctx) error {
	tools, err := utils.GetOrSetCache("web3tools", 24*time.Hour, func() ([]Web3Tools, error) {
		return GetAllWeb3Tools()
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

func GetWeb3ToolStatsHandler(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("web3tools_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return GetWeb3ToolStats()
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

func GetWeb3ToolsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	tool, err := GetWeb3ToolsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Web3Tools not found",
		})
	}

	return c.JSON(tool)
}

func InsertWeb3ToolsHandler(c *fiber.Ctx) error {
	var req Web3Tools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := InsertWeb3Tools(
		req.Name,
		req.Description,
		req.Category,
		req.Chains,
		req.ImageURL,
		req.Website,
		req.Media,
		req.AddedBy,
		req.Twitter,
		req.Instagram,
		req.Discord,
		req.Telegram,
		req.Youtube,
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

func UpdateWeb3ToolsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req Web3Tools

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := Web3Tools{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Chains:      req.Chains,
		ImageURL:    req.ImageURL,
		Website:     req.Website,
		Media:       req.Media,
		AddedBy:     req.AddedBy,
		Twitter:     req.Twitter,
		Instagram:   req.Instagram,
		Discord:     req.Discord,
		Telegram:    req.Telegram,
		Youtube:     req.Youtube,
	}

	updatedTool, err := UpdateWeb3ToolsByID(id, updateData)
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

func DeleteWeb3ToolsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = DeleteWeb3ToolsByID(id)
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