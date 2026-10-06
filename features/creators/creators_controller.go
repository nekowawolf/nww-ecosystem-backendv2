package creators

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateCreatorsCache() {
	utils.InvalidateCache("creators", "creators_stats")
}

func GetAllCreatorsHandler(c *fiber.Ctx) error {
	creators, err := utils.GetOrSetCache("creators", 24*time.Hour, func() ([]Creators, error) {
		return GetAllCreators()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    creators,
	})
}

func GetCreatorsStatsHandler(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("creators_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return GetCreatorsStats()
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

func GetCreatorsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	creatorItem, err := GetCreatorsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Creator not found",
		})
	}

	return c.JSON(creatorItem)
}

func InsertCreatorsHandler(c *fiber.Ctx) error {
	var req Creators

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := InsertCreators(
		req.Name,
		req.Description,
		req.ImageURL,
		req.Website,
		req.Category,
		req.Language,
		req.OpenToWork,
		req.Socials,
		req.Platforms,
		req.AddedBy,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert Creator",
		})
	}

	invalidateCreatorsCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Creator created successfully",
		"insertedID": insertedID,
	})
}

func UpdateCreatorsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req Creators

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := Creators{
		Name:          req.Name,
		Description:   req.Description,
		ImageURL:      req.ImageURL,
		Website:       req.Website,
		Category:      req.Category,
		Language:      req.Language,
		OpenToWork:    req.OpenToWork,
		Socials:       req.Socials,
		Platforms:     req.Platforms,
		AddedBy:       req.AddedBy,
	}

	updatedCreator, err := UpdateCreatorsByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Creator not found or could not be updated",
		})
	}

	invalidateCreatorsCache()
	return c.JSON(fiber.Map{
		"message": "Creator updated successfully",
		"data":    updatedCreator,
	})
}

func DeleteCreatorsByIDHandler(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = DeleteCreatorsByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateCreatorsCache()
	return c.JSON(fiber.Map{
		"message": "Creator deleted successfully",
	})
}