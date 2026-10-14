package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateNetBlogCache() {
	utils.InvalidateCache("netblog")
	utils.InvalidateCachePrefix("netblog_")
}

func netBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NetBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A NetBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process NetBlog"})
}

func GetAllNetBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("netblog", 24*time.Hour, GetAllNetBlogs)
	if err != nil {
		return netBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetNetBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !netBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("netblog_slug_"+slug, 24*time.Hour, func() (*NetBlog, error) {
		return GetNetBlogBySlug(slug)
	})
	if err != nil {
		return netBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetNetBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("netblog_id_"+id.Hex(), 24*time.Hour, func() (*NetBlog, error) {
		return GetNetBlogByID(id)
	})
	if err != nil {
		return netBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertNetBlogHandler(c *fiber.Ctx) error {
	var req netBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateNetBlog(&req.NetBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertNetBlog(req.NetBlog)
	if err != nil {
		return netBlogError(c, err)
	}
	invalidateNetBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "NetBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateNetBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req netBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateNetBlog(&req.NetBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateNetBlogByID(id, req.NetBlog)
	if err != nil {
		return netBlogError(c, err)
	}
	invalidateNetBlogCache()
	return c.JSON(fiber.Map{"message": "NetBlog updated successfully", "data": blog})
}

func DeleteNetBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteNetBlogByID(id); err != nil {
		return netBlogError(c, err)
	}
	invalidateNetBlogCache()
	return c.JSON(fiber.Map{"message": "NetBlog deleted successfully"})
}
