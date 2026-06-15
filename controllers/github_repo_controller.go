package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func invalidateGithubRepoCache() {
	utils.InvalidateCache("githubrepo", "githubrepo_stats")
}

func GetAllGithubRepos(c *fiber.Ctx) error {
	repos, err := utils.GetOrSetCache("githubrepo", 24*time.Hour, func() ([]models.GithubRepo, error) {
		return module.GetAllGithubRepos()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    repos,
	})
}

func GetGithubRepoStats(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("githubrepo_stats", 24*time.Hour, func() (map[string]interface{}, error) {
		return module.GetGithubRepoStats()
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

func GetGithubRepoByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	repo, err := module.GetGithubRepoByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "GithubRepo not found",
		})
	}

	return c.JSON(repo)
}

func InsertGithubRepo(c *fiber.Ctx) error {
	var req models.GithubRepo

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := module.InsertGithubRepo(
		req.Name,
		req.Description,
		req.Category,
		req.RepoURL,
		req.Owner,
		req.RepoName,
		req.Twitter,
		req.Discord,
		req.Telegram,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert GithubRepo",
		})
	}

	invalidateGithubRepoCache()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "GithubRepo created successfully",
		"insertedID": insertedID,
	})
}

func UpdateGithubRepoByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req models.GithubRepo

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := models.GithubRepo{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		RepoURL:     req.RepoURL,
		Owner:       req.Owner,
		RepoName:    req.RepoName,
		Twitter:     req.Twitter,
		Discord:     req.Discord,
		Telegram:    req.Telegram,
	}

	updatedRepo, err := module.UpdateGithubRepoByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "GithubRepo not found or could not be updated",
		})
	}

	invalidateGithubRepoCache()
	return c.JSON(fiber.Map{
		"message": "GithubRepo updated successfully",
		"data":    updatedRepo,
	})
}

func DeleteGithubRepoByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = module.DeleteGithubRepoByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateGithubRepoCache()
	return c.JSON(fiber.Map{
		"message": "GithubRepo deleted successfully",
	})
}