package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateAirdropBlogCache() {
	utils.InvalidateCache("airdropblog")
	utils.InvalidateCachePrefix("airdropblog_")
}

func airdropBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "AirdropBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A AirdropBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process AirdropBlog"})
}

func GetAllAirdropBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("airdropblog", 24*time.Hour, GetAllAirdropBlogs)
	if err != nil {
		return airdropBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetAirdropBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !airdropBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("airdropblog_slug_"+slug, 24*time.Hour, func() (*AirdropBlog, error) {
		return GetAirdropBlogBySlug(slug)
	})
	if err != nil {
		return airdropBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetAirdropBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("airdropblog_id_"+id.Hex(), 24*time.Hour, func() (*AirdropBlog, error) {
		return GetAirdropBlogByID(id)
	})
	if err != nil {
		return airdropBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertAirdropBlogHandler(c *fiber.Ctx) error {
	var req airdropBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateAirdropBlog(&req.AirdropBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertAirdropBlog(req.AirdropBlog)
	if err != nil {
		return airdropBlogError(c, err)
	}
	invalidateAirdropBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "AirdropBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateAirdropBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req airdropBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateAirdropBlog(&req.AirdropBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateAirdropBlogByID(id, req.AirdropBlog)
	if err != nil {
		return airdropBlogError(c, err)
	}
	invalidateAirdropBlogCache()
	return c.JSON(fiber.Map{"message": "AirdropBlog updated successfully", "data": blog})
}

func DeleteAirdropBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteAirdropBlogByID(id); err != nil {
		return airdropBlogError(c, err)
	}
	invalidateAirdropBlogCache()
	return c.JSON(fiber.Map{"message": "AirdropBlog deleted successfully"})
}