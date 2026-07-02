package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func GetAllNotes(c *fiber.Ctx) error {
	notes, err := module.GetAllNotes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    notes,
	})
}

func GetNoteByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	note, err := module.GetNoteByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Note not found",
		})
	}

	return c.JSON(note)
}

func InsertNote(c *fiber.Ctx) error {
	var req models.Note

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	insertedID := module.InsertNote(
		req.Title,
		req.Content,
		req.Type,
	)

	if insertedID == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert Note",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Note created successfully",
		"insertedID": insertedID,
	})
}

func UpdateNoteByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	var req models.Note

	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}

	updateData := models.Note{
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
	}

	updatedNote, err := module.UpdateNoteByID(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Note not found or could not be updated",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Note updated successfully",
		"data":    updatedNote,
	})
}

func DeleteNoteByID(c *fiber.Ctx) error {
	id, err := utils.ParseObjectID(c, "id")
	if err != nil {
		return err
	}

	err = module.DeleteNoteByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Note deleted successfully",
	})
}
