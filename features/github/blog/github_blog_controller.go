package blog

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func invalidateGithubBlogCache() {
	utils.InvalidateCache("githubblog")
	utils.InvalidateCachePrefix("githubblog_")
}

func githubBlogError(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "GithubBlog not found"})
	}
	if mongo.IsDuplicateKeyError(err) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "A GithubBlog with this slug already exists"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process GithubBlog"})
}

func GetAllGithubBlogsHandler(c *fiber.Ctx) error {
	blogs, err := utils.GetOrSetCache("githubblog", 24*time.Hour, GetAllGithubBlogs)
	if err != nil {
		return githubBlogError(c, err)
	}
	return c.JSON(fiber.Map{"message": "Data retrieved successfully", "data": blogs})
}

func GetGithubBlogBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if len(slug) > 200 || !githubBlogSlugPattern.MatchString(slug) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid slug format"})
	}
	blog, err := utils.GetOrSetCache("githubblog_slug_"+slug, 24*time.Hour, func() (*GithubBlog, error) {
		return GetGithubBlogBySlug(slug)
	})
	if err != nil {
		return githubBlogError(c, err)
	}
	return c.JSON(blog)
}

func GetGithubBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	blog, err := utils.GetOrSetCache("githubblog_id_"+id.Hex(), 24*time.Hour, func() (*GithubBlog, error) {
		return GetGithubBlogByID(id)
	})
	if err != nil {
		return githubBlogError(c, err)
	}
	return c.JSON(blog)
}

func InsertGithubBlogHandler(c *fiber.Ctx) error {
	var req githubBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateGithubBlog(&req.GithubBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := InsertGithubBlog(req.GithubBlog)
	if err != nil {
		return githubBlogError(c, err)
	}
	invalidateGithubBlogCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "GithubBlog created successfully", "insertedID": blog.ID, "data": blog,
	})
}

func UpdateGithubBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var req githubBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := validateGithubBlog(&req.GithubBlog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	blog, err := UpdateGithubBlogByID(id, req.GithubBlog)
	if err != nil {
		return githubBlogError(c, err)
	}
	invalidateGithubBlogCache()
	return c.JSON(fiber.Map{"message": "GithubBlog updated successfully", "data": blog})
}

func DeleteGithubBlogByIDHandler(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := DeleteGithubBlogByID(id); err != nil {
		return githubBlogError(c, err)
	}
	invalidateGithubBlogCache()
	return c.JSON(fiber.Map{"message": "GithubBlog deleted successfully"})
}