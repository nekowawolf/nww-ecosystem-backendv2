package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateGuildBlogCache() {
	utils.InvalidateCache("guildblog")
	utils.InvalidateCachePrefix("guildblog_")
}

func guildBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "GuildBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A GuildBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process GuildBlog"})
}

func GetAllGuildBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("guildblog", 24*time.Hour, GetAllGuildBlogs)
	if err != nil {
		return guildBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetGuildBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !guildBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("guildblog_slug_"+slug, 24*time.Hour, func() (*GuildBlog, error) {
		return GetGuildBlogBySlug(slug)
	})
	if err != nil {
		return guildBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetGuildBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("guildblog_id_"+id.Hex(), 24*time.Hour, func() (*GuildBlog, error) {
		return GetGuildBlogByID(id)
	})
	if err != nil {
		return guildBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertGuildBlogHandler(c *fiber.Ctx) error {
	var req guildBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateGuildBlog(&req.GuildBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertGuildBlog(req.GuildBlog)
	if err != nil {
		return guildBlogError(c, err)
	}
	invalidateGuildBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "GuildBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateGuildBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req guildBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateGuildBlog(&req.GuildBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateGuildBlogByID(id, req.GuildBlog)
	if err != nil {
		return guildBlogError(c, err)
	}
	invalidateGuildBlogCache()
	return c.JSON(fiber.Map{"message": "GuildBlog updated successfully", "data": blog})
}

func DeleteGuildBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteGuildBlogByID(id); err != nil {
		return guildBlogError(c, err)
	}
	invalidateGuildBlogCache()
	return c.JSON(fiber.Map{"message": "GuildBlog deleted successfully"})
}
