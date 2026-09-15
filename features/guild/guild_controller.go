package guild

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateGuildCache() {
	utils.InvalidateCache("guild", "guild_stats")
}

func GetAllGuildHandler(c *fiber.Ctx) error {
	guilds, err := utils.GetOrSetCache("guild", 24*time.Hour, func() ([]Guild, error) {
		return GetAllGuild()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    guilds,
	})
}

func GetGuildStatsHandler(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("guild_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return GetGuildStats()
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

func GetGuildByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	guild, err := GetGuildByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Guild not found",
		})
	}

	return c.JSON(guild)
}

func InsertGuildHandler(c *fiber.Ctx) error {
	var req Guild

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := InsertGuild(req)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert Guild",
		})
	}

	invalidateGuildCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Guild created successfully",
		"insertedID": insertedID,
	})
}

func UpdateGuildByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req Guild

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := Guild{
		Name:        req.Name,
		Description: req.Description,
		Platform:    req.Platform,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Website:     req.Website,
		Link:        req.Link,
		Socials:     req.Socials,
	}

	updatedGuild, err := UpdateGuildByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Guild not found or could not be updated",
		})
	}

	invalidateGuildCache()
	return c.JSON(fiber.Map{
		"message": "Guild updated successfully",
		"data":    updatedGuild,
	})
}

func DeleteGuildByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = DeleteGuildByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateGuildCache()
	return c.JSON(fiber.Map{
		"message": "Guild deleted successfully",
	})
}