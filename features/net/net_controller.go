package net

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateNetCache() {
	utils.InvalidateCache("net", "net_stats")
}

func GetAllNetHandler(c *fiber.Ctx) error {
	nets, err := utils.GetOrSetCache("net", 24*time.Hour, func() ([]Net, error) {
		return GetAllNet()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    nets,
	})
}

func GetNetStatsHandler(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("net_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return GetNetStats()
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

func GetNetByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	netItem, err := GetNetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Net not found",
		})
	}

	return c.JSON(netItem)
}

func InsertNetHandler(c *fiber.Ctx) error {
	var req Net

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := InsertNet(
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
			"error": "Failed to insert Net",
		})
	}

	invalidateNetCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Net created successfully",
		"insertedID": insertedID,
	})
}

func UpdateNetByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req Net

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := Net{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Website:     req.Website,
		Categories:  req.Categories,
		Media:       req.Media,
		Socials:     req.Socials,
		AddedBy:     req.AddedBy,
	}

	updatedNet, err := UpdateNetByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Net not found or could not be updated",
		})
	}

	invalidateNetCache()
	return c.JSON(fiber.Map{
		"message": "Net updated successfully",
		"data":    updatedNet,
	})
}

func DeleteNetByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = DeleteNetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateNetCache()
	return c.JSON(fiber.Map{
		"message": "Net deleted successfully",
	})
}