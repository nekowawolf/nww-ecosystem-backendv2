package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateCreatorsBlogCache() {
	utils.InvalidateCache("creatorsblog")
	utils.InvalidateCachePrefix("creatorsblog_")
}

func creatorsBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "CreatorsBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A CreatorsBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process CreatorsBlog"})
}

func GetAllCreatorsBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("creatorsblog", 24*time.Hour, GetAllCreatorsBlogs)
	if err != nil {
		return creatorsBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetCreatorsBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !creatorsBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("creatorsblog_slug_"+slug, 24*time.Hour, func() (*CreatorsBlog, error) {
		return GetCreatorsBlogBySlug(slug)
	})
	if err != nil {
		return creatorsBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetCreatorsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("creatorsblog_id_"+id.Hex(), 24*time.Hour, func() (*CreatorsBlog, error) {
		return GetCreatorsBlogByID(id)
	})
	if err != nil {
		return creatorsBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertCreatorsBlogHandler(c *fiber.Ctx) error {
	var req creatorsBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateCreatorsBlog(&req.CreatorsBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertCreatorsBlog(req.CreatorsBlog)
	if err != nil {
		return creatorsBlogError(c, err)
	}
	invalidateCreatorsBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "CreatorsBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateCreatorsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req creatorsBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateCreatorsBlog(&req.CreatorsBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateCreatorsBlogByID(id, req.CreatorsBlog)
	if err != nil {
		return creatorsBlogError(c, err)
	}
	invalidateCreatorsBlogCache()
	return c.JSON(fiber.Map{"message": "CreatorsBlog updated successfully", "data": blog})
}

func DeleteCreatorsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteCreatorsBlogByID(id); err != nil {
		return creatorsBlogError(c, err)
	}
	invalidateCreatorsBlogCache()
	return c.JSON(fiber.Map{"message": "CreatorsBlog deleted successfully"})
}
