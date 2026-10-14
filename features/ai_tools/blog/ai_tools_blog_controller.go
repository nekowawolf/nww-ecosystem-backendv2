package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateAIToolsBlogCache() {
	utils.InvalidateCache("aitoolsblog")
	utils.InvalidateCachePrefix("aitoolsblog_")
}

func aiToolsBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "AIToolsBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A AIToolsBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process AIToolsBlog"})
}

func GetAllAIToolsBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("aitoolsblog", 24*time.Hour, GetAllAIToolsBlogs)
	if err != nil {
		return aiToolsBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetAIToolsBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !aiToolsBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("aitoolsblog_slug_"+slug, 24*time.Hour, func() (*AIToolsBlog, error) {
		return GetAIToolsBlogBySlug(slug)
	})
	if err != nil {
		return aiToolsBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetAIToolsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("aitoolsblog_id_"+id.Hex(), 24*time.Hour, func() (*AIToolsBlog, error) {
		return GetAIToolsBlogByID(id)
	})
	if err != nil {
		return aiToolsBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertAIToolsBlogHandler(c *fiber.Ctx) error {
	var req aiToolsBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateAIToolsBlog(&req.AIToolsBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertAIToolsBlog(req.AIToolsBlog)
	if err != nil {
		return aiToolsBlogError(c, err)
	}
	invalidateAIToolsBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "AIToolsBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateAIToolsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req aiToolsBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateAIToolsBlog(&req.AIToolsBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateAIToolsBlogByID(id, req.AIToolsBlog)
	if err != nil {
		return aiToolsBlogError(c, err)
	}
	invalidateAIToolsBlogCache()
	return c.JSON(fiber.Map{"message": "AIToolsBlog updated successfully", "data": blog})
}

func DeleteAIToolsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteAIToolsBlogByID(id); err != nil {
		return aiToolsBlogError(c, err)
	}
	invalidateAIToolsBlogCache()
	return c.JSON(fiber.Map{"message": "AIToolsBlog deleted successfully"})
}