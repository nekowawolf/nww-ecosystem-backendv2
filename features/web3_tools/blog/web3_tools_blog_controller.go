package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateWeb3ToolsBlogCache() {
	utils.InvalidateCache("web3toolsblog")
	utils.InvalidateCachePrefix("web3toolsblog_")
}

func web3ToolsBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Web3ToolsBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A Web3ToolsBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process Web3ToolsBlog"})
}

func GetAllWeb3ToolsBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("web3toolsblog", 24*time.Hour, GetAllWeb3ToolsBlogs)
	if err != nil {
		return web3ToolsBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetWeb3ToolsBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !web3ToolsBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("web3toolsblog_slug_"+slug, 24*time.Hour, func() (*Web3ToolsBlog, error) {
		return GetWeb3ToolsBlogBySlug(slug)
	})
	if err != nil {
		return web3ToolsBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetWeb3ToolsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("web3toolsblog_id_"+id.Hex(), 24*time.Hour, func() (*Web3ToolsBlog, error) {
		return GetWeb3ToolsBlogByID(id)
	})
	if err != nil {
		return web3ToolsBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertWeb3ToolsBlogHandler(c *fiber.Ctx) error {
	var req web3ToolsBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateWeb3ToolsBlog(&req.Web3ToolsBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertWeb3ToolsBlog(req.Web3ToolsBlog)
	if err != nil {
		return web3ToolsBlogError(c, err)
	}
	invalidateWeb3ToolsBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Web3ToolsBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateWeb3ToolsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req web3ToolsBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateWeb3ToolsBlog(&req.Web3ToolsBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateWeb3ToolsBlogByID(id, req.Web3ToolsBlog)
	if err != nil {
		return web3ToolsBlogError(c, err)
	}
	invalidateWeb3ToolsBlogCache()
	return c.JSON(fiber.Map{"message": "Web3ToolsBlog updated successfully", "data": blog})
}

func DeleteWeb3ToolsBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteWeb3ToolsBlogByID(id); err != nil {
		return web3ToolsBlogError(c, err)
	}
	invalidateWeb3ToolsBlogCache()
	return c.JSON(fiber.Map{"message": "Web3ToolsBlog deleted successfully"})
}
