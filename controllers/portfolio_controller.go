package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func GetPortfolio(c *fiber.Ctx) error {
	data, err := utils.GetOrSetCache("portfolio", 24*time.Hour, func() (*models.Portfolio, error) {
		return module.GetPortfolio()
	})
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Portfolio not found"})
	}
	return c.JSON(data)
}

func UpdatePortfolio(c *fiber.Ctx) error {
	var req models.Portfolio
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.UpdatePortfolio(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Portfolio updated"})
}

func UpdateHeroProfile(c *fiber.Ctx) error {
	var req models.HeroProfile
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.UpdateHeroProfile(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update hero profile failed"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Hero profile updated"})
}

func AddCertificate(c *fiber.Ctx) error {
	var req models.Certificate
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.AddCertificate(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add certificate"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Certificate added"})
}

func AddDesign(c *fiber.Ctx) error {
	var req models.Design
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.AddDesign(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add design"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Design added"})
}

func AddProject(c *fiber.Ctx) error {
	var req models.Project
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.AddProject(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add project"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Project added"})
}

func AddExperience(c *fiber.Ctx) error {
	var req models.Experience
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.AddExperience(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add experience"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Experience added"})
}

func AddEducation(c *fiber.Ctx) error {
	var req models.Education
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.AddEducation(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add education"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Education added"})
}

func AddTechSkill(c *fiber.Ctx) error {
	var req models.SkillItem
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.AddTechSkill(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add tech skill"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Tech skill added"})
}

func DeleteCertificate(c *fiber.Ctx) error {
	if err := module.DeleteCertificate(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete certificate"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Certificate deleted"})
}

func DeleteDesign(c *fiber.Ctx) error {
	if err := module.DeleteDesign(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete design"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Design deleted"})
}

func DeleteProject(c *fiber.Ctx) error {
	if err := module.DeleteProject(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete project"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Project deleted"})
}

func DeleteExperience(c *fiber.Ctx) error {
	if err := module.DeleteExperience(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete experience"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Experience deleted"})
}

func DeleteEducation(c *fiber.Ctx) error {
	if err := module.DeleteEducation(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete education"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Education deleted"})
}

func AddDesignSkill(c *fiber.Ctx) error {
	var req models.SkillItem
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := module.AddDesignSkill(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add design skill"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Design skill added"})
}

func DeleteDesignSkill(c *fiber.Ctx) error {
	if err := module.DeleteDesignSkill(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete design skill"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Design skill deleted"})
}

func DeleteTechSkill(c *fiber.Ctx) error {
	if err := module.DeleteTechSkill(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete tech skill"})
	}
	utils.InvalidateCache("portfolio")
	return c.JSON(fiber.Map{"message": "Tech skill deleted"})
}