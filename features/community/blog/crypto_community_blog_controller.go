package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateCryptoCommunityBlogCache() {
	utils.InvalidateCache("cryptocommunityblog")
	utils.InvalidateCachePrefix("cryptocommunityblog_")
}

func cryptoCommunityBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "CryptoCommunityBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A CryptoCommunityBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process CryptoCommunityBlog"})
}

func GetAllCryptoCommunityBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("cryptocommunityblog", 24*time.Hour, GetAllCryptoCommunityBlogs)
	if err != nil {
		return cryptoCommunityBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetCryptoCommunityBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !cryptoCommunityBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("cryptocommunityblog_slug_"+slug, 24*time.Hour, func() (*CryptoCommunityBlog, error) {
		return GetCryptoCommunityBlogBySlug(slug)
	})
	if err != nil {
		return cryptoCommunityBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetCryptoCommunityBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("cryptocommunityblog_id_"+id.Hex(), 24*time.Hour, func() (*CryptoCommunityBlog, error) {
		return GetCryptoCommunityBlogByID(id)
	})
	if err != nil {
		return cryptoCommunityBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertCryptoCommunityBlogHandler(c *fiber.Ctx) error {
	var req cryptoCommunityBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateCryptoCommunityBlog(&req.CryptoCommunityBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertCryptoCommunityBlog(req.CryptoCommunityBlog)
	if err != nil {
		return cryptoCommunityBlogError(c, err)
	}
	invalidateCryptoCommunityBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "CryptoCommunityBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateCryptoCommunityBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req cryptoCommunityBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateCryptoCommunityBlog(&req.CryptoCommunityBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateCryptoCommunityBlogByID(id, req.CryptoCommunityBlog)
	if err != nil {
		return cryptoCommunityBlogError(c, err)
	}
	invalidateCryptoCommunityBlogCache()
	return c.JSON(fiber.Map{"message": "CryptoCommunityBlog updated successfully", "data": blog})
}

func DeleteCryptoCommunityBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteCryptoCommunityBlogByID(id); err != nil {
		return cryptoCommunityBlogError(c, err)
	}
	invalidateCryptoCommunityBlogCache()
	return c.JSON(fiber.Map{"message": "CryptoCommunityBlog deleted successfully"})
}
