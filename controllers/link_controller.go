package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ==================== PROFILE CONTROLLERS ====================

func GetProfile(c *fiber.Ctx) error {
	profile, err := utils.GetOrSetCache("profilelink", 24*time.Hour, func() (*models.Profile, error) {
		return module.GetProfile()
	})
	if err != nil {
		// Return default profile if not found
		return c.JSON(fiber.Map{
			"name":       "nekowawolf",
			"username":   "nekowawolf",
			"bio":        "Professional Coder (vibe coding)",
			"avatar_url": "https://nekowawolf.github.io/cdn-images/images/2025/1763530019_113094795.jpeg",
			"cover_url":  "https://nekowawolf.github.io/cdn-images/images/2026/1775599464_bg_link.png",
			"links": fiber.Map{
				"github":    "https://github.com/nekowawolf",
				"twitter":   "https://x.com/nekowawolf_",
				"tiktok":    "https://tiktok.com/@nekowawolf",
				"website":   "https://nekowawolf.xyz/",
				"instagram": "https://instagram.com/nekowawolf",
			},
		})
	}

	return c.JSON(fiber.Map{
		"name":       profile.Name,
		"username":   profile.Username,
		"bio":        profile.Bio,
		"avatar_url": profile.AvatarURL,
		"cover_url":  profile.CoverURL,
		"links":      profile.Links,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	var req models.Profile
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := module.UpdateProfile(req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.InvalidateCache("profilelink")
	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

func GetPostStats(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("postslink_stats", 24*time.Hour, func() (interface{}, error) {
		return module.GetPostStats()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve stats",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Stats retrieved successfully",
		"data":    stats,
	})
}

// ==================== POSTS CONTROLLERS ====================

func GetAllPosts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 5)
	category := c.Query("category", "")
	search := c.Query("search", "")

	cacheKey := fmt.Sprintf("postslink:%d:%d:%s:%s", page, limit, category, search)

	posts, err := utils.GetOrSetCache(cacheKey, 24*time.Hour, func() (interface{}, error) {
		return module.GetPostsPaginated(page, limit, category, search)
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve posts",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Posts retrieved successfully",
		"data":    posts,
	})
}

func GetPostByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	post, err := module.GetPostByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	sessionID := c.Get("X-Session-ID")
	if sessionID == "" {
		sessionID = c.IP() + c.Get("User-Agent")
	}

	_ = module.IncrementPostView(id, sessionID)

	viewCount, _ := module.GetPostViewCount(id)
	post.Views = int(viewCount)

	return c.JSON(fiber.Map{
		"message": "Post retrieved successfully",
		"data":    post,
	})
}

func CreatePost(c *fiber.Ctx) error {
	var req models.LinkPost
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Caption == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Caption is required",
		})
	}

	if req.Category == "" {
		req.Category = "all"
	}

	insertedID, err := module.InsertPost(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if objectID, ok := insertedID.(primitive.ObjectID); ok {
		utils.InvalidateCache("postslink_stats")
		utils.InvalidateCachePrefix("postslink:")
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":     "Post created successfully",
			"inserted_id": objectID.Hex(),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Failed to retrieve inserted ID",
	})
}

func UpdatePost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var req models.LinkPost
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := module.UpdatePost(id, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.InvalidateCache("postslink_stats")
	utils.InvalidateCachePrefix("postslink:")
	return c.JSON(fiber.Map{
		"message": "Post updated successfully",
	})
}

func DeletePost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	if err := module.DeletePost(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.InvalidateCache("postslink_stats")
	utils.InvalidateCachePrefix("postslink:")
	return c.JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}