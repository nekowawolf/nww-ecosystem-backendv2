package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func invalidateAirdropCache() {
	utils.InvalidateCache("freeairdrop", "paidairdrop", "allairdrop", "allairdrop_stats")
}

func GetAllAirdropHandler(c *fiber.Ctx) error {
	data, err := utils.GetOrSetCache("allairdrop", 24*time.Hour, func() ([]interface{}, error) {
		return module.GetAllAirdrop()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve all Airdrop data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAllAirdropStatsHandler(c *fiber.Ctx) error {
	stats, err := utils.GetOrSetCache("allairdrop_stats", 24*time.Hour, func() (map[string]int, error) {
		return module.GetAllAirdropStats()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stats retrieved successfully",
		"data":    stats,
	})
}

func GetAirdropFreeHandler(c *fiber.Ctx) error {
	data, err := utils.GetOrSetCache("freeairdrop", 24*time.Hour, func() ([]models.AirdropFree, error) {
		return module.GetAllAirdropFree()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve AirdropFree data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAirdropPaidHandler(c *fiber.Ctx) error {
	data, err := utils.GetOrSetCache("paidairdrop", 24*time.Hour, func() ([]models.AirdropPaid, error) {
		return module.GetAllAirdropPaid()
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve AirdropFree data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAllAirdropByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	data, err := module.GetAllAirdropByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAirdropFreeByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	data, err := module.GetAirdropFreeByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve AirdropFree by ID",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAirdropPaidByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	data, err := module.GetAirdropPaidByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve AirdropPaid by ID",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAllAirdropByNameHandler(c *fiber.Ctx) error {
	name := c.Params("name")

	data, err := module.GetAllAirdropByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve Airdrop by Name",
		})
	}

	if len(data) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Airdrop found with the specified name",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAirdropFreeByNameHandler(c *fiber.Ctx) error {
	name := c.Params("name")

	data, err := module.GetAirdropFreeByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve AirdropFree by Name",
		})
	}

	if len(data) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No AirdropFree found with the specified name",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func GetAirdropPaidByNameHandler(c *fiber.Ctx) error {
	name := c.Params("name")

	data, err := module.GetAirdropPaidByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve AirdropPaid by Name",
		})
	}

	if len(data) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No AirdropPaid found with the specified name",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data retrieved successfully",
		"data":    data,
	})
}

func InsertAirdropFreeHandler(c *fiber.Ctx) error {
	var reqAirdrop models.AirdropFree

	if err := c.BodyParser(&reqAirdrop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if reqAirdrop.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	insertedID, err := module.InsertAirdropFree(
		reqAirdrop.Name,
		reqAirdrop.Task,
		reqAirdrop.Link,
		reqAirdrop.Level,
		reqAirdrop.Status,
		reqAirdrop.Backed,
		reqAirdrop.Funds,
		reqAirdrop.Supply,
		reqAirdrop.Fdv,
		reqAirdrop.MarketCap,
		reqAirdrop.Vesting,
		reqAirdrop.LinkClaim,
		reqAirdrop.LinkDiscord,
		reqAirdrop.LinkTwitter,
		reqAirdrop.LinkTelegram,
		reqAirdrop.ImageURL,
		reqAirdrop.Description,
		reqAirdrop.LinkGuide,
		reqAirdrop.Price,
		reqAirdrop.USDIncome,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert AirdropFree",
		})
	}

	if objectID, ok := insertedID.(primitive.ObjectID); ok {
		invalidateAirdropCache()
		return c.JSON(fiber.Map{
			"message":     "AirdropFree inserted successfully",
			"inserted_id": objectID.Hex(),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Failed to retrieve inserted ID",
	})
}

func InsertAirdropPaidHandler(c *fiber.Ctx) error {
	var reqAirdrop models.AirdropPaid

	if err := c.BodyParser(&reqAirdrop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if reqAirdrop.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	insertedID, err := module.InsertAirdropPaid(
		reqAirdrop.Name,
		reqAirdrop.Task,
		reqAirdrop.Link,
		reqAirdrop.Level,
		reqAirdrop.Status,
		reqAirdrop.Backed,
		reqAirdrop.Funds,
		reqAirdrop.Supply,
		reqAirdrop.Fdv,
		reqAirdrop.MarketCap,
		reqAirdrop.Vesting,
		reqAirdrop.LinkClaim,
		reqAirdrop.LinkDiscord,
		reqAirdrop.LinkTwitter,
		reqAirdrop.LinkTelegram,
		reqAirdrop.ImageURL,
		reqAirdrop.Description,
		reqAirdrop.LinkGuide,
		reqAirdrop.Price,
		reqAirdrop.USDIncome,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert AirdropPaid",
		})
	}

	if objectID, ok := insertedID.(primitive.ObjectID); ok {
		invalidateAirdropCache()
		return c.JSON(fiber.Map{
			"message":     "AirdropPaid inserted successfully",
			"inserted_id": objectID.Hex(),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Failed to retrieve inserted ID",
	})
}

func UpdateAllAirdropByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var updateData struct {
		Name        string  `json:"name"`
		Task        string  `json:"task"`
		Link        string  `json:"link"`
		Level       string  `json:"level"`
		Status      string  `json:"status"`
		Backed      string  `json:"backed"`
		Funds       string  `json:"funds"`
		Supply      string  `json:"supply"`
		Fdv         string  `json:"fdv"`
		MarketCap   string  `json:"market_cap"`
		Vesting     string  `json:"vesting"`
		LinkClaim   string  `json:"link_claim"`
		LinkDiscord string  `json:"link_discord"`
		LinkTwitter string  `json:"link_twitter"`
		LinkTelegram string `json:"link_telegram"`
		ImageURL    string  `json:"image_url"`
		Description string  `json:"description"`
		LinkGuide   string  `json:"link_guide"`
		Price       float64 `json:"price"`
		USDIncome   int     `json:"usd_income"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	currentAirdrop, err := module.GetAllAirdropByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var finalStatus string
	if updateData.Status == "" {
		if freeAirdrop, ok := currentAirdrop.(models.AirdropFree); ok {
			finalStatus = freeAirdrop.Status
		} else if paidAirdrop, ok := currentAirdrop.(models.AirdropPaid); ok {
			finalStatus = paidAirdrop.Status
		} else {
			finalStatus = updateData.Status
		}
	} else {
		finalStatus = updateData.Status
	}

	err = module.UpdateAllAirdropByID(
		id,
		updateData.Name,
		updateData.Task,
		updateData.Link,
		updateData.Level,
		finalStatus,
		updateData.Backed,
		updateData.Funds,
		updateData.Supply,
		updateData.Fdv,
		updateData.MarketCap,
		updateData.Vesting,
		updateData.LinkClaim,
		updateData.LinkDiscord,
		updateData.LinkTwitter,
		updateData.LinkTelegram,
		updateData.ImageURL,
		updateData.Description,
		updateData.LinkGuide,
		updateData.Price,
		updateData.USDIncome,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateAirdropCache()
	return c.JSON(fiber.Map{
		"message": "Airdrop updated successfully",
	})
}

func UpdateAirdropFreeByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var updateData struct {
		Name        string  `json:"name"`
		Task        string  `json:"task"`
		Link        string  `json:"link"`
		Level       string  `json:"level"`
		Status      string  `json:"status"`
		Backed      string  `json:"backed"`
		Funds       string  `json:"funds"`
		Supply      string  `json:"Supply"`
		Fdv         string  `json:"fdv"`
		MarketCap   string  `json:"market_cap"`
		Vesting     string  `json:"vesting"`
		LinkClaim   string  `json:"link_claim"`
		LinkDiscord string  `json:"link_discord"`
		LinkTwitter string  `json:"link_twitter"`
		LinkTelegram string `json:"link_telegram"`
		ImageURL    string  `json:"image_url"`
		Description string  `json:"description"`
		LinkGuide   string  `json:"link_guide"`
		Price       float64 `json:"price"`
		USDIncome   int     `json:"usd_income"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	err = module.UpdateAirdropFreeByID(
		id,
		updateData.Name,
		updateData.Task,
		updateData.Link,
		updateData.Level,
		updateData.Status,
		updateData.Backed,
		updateData.Funds,
		updateData.Supply,
		updateData.Fdv,
		updateData.MarketCap,
		updateData.Vesting,
		updateData.LinkClaim,
		updateData.LinkDiscord,
		updateData.LinkTwitter,
		updateData.LinkTelegram,
		updateData.ImageURL,
		updateData.Description,
		updateData.LinkGuide,
		updateData.Price,
		updateData.USDIncome,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update AirdropFree by ID",
		})
	}

	invalidateAirdropCache()
	return c.JSON(fiber.Map{
		"message": "AirdropFree updated successfully",
	})
}

func UpdateAirdropPaidByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var updateData struct {
		Name        string  `json:"name"`
		Task        string  `json:"task"`
		Link        string  `json:"link"`
		Level       string  `json:"level"`
		Status      string  `json:"status"`
		Backed      string  `json:"backed"`
		Funds       string  `json:"funds"`
		Supply      string  `json:"Supply"`
		Fdv         string  `json:"fdv"`
		MarketCap   string  `json:"market_cap"`
		Vesting     string  `json:"vesting"`
		LinkClaim   string  `json:"link_claim"`
		LinkDiscord string  `json:"link_discord"`
		LinkTwitter string  `json:"link_twitter"`
		LinkTelegram string `json:"link_telegram"`
		ImageURL    string  `json:"image_url"`
		Description string  `json:"description"`
		LinkGuide   string  `json:"link_guide"`
		Price       float64 `json:"price"`
		USDIncome   int     `json:"usd_income"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	err = module.UpdateAirdropPaidByID(
		id,
		updateData.Name,
		updateData.Task,
		updateData.Link,
		updateData.Level,
		updateData.Status,
		updateData.Backed,
		updateData.Funds,
		updateData.Supply,
		updateData.Fdv,
		updateData.MarketCap,
		updateData.Vesting,
		updateData.LinkClaim,
		updateData.LinkDiscord,
		updateData.LinkTwitter,
		updateData.LinkTelegram,
		updateData.ImageURL,
		updateData.Description,
		updateData.LinkGuide,
		updateData.Price,
		updateData.USDIncome,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update AirdropPaid by ID",
		})
	}

	invalidateAirdropCache()
	return c.JSON(fiber.Map{
		"message": "AirdropPaid updated successfully",
	})
}

func DeleteAllAirdropByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	err = module.DeleteAllAirdropByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	invalidateAirdropCache()
	return c.JSON(fiber.Map{
		"message": "Airdrop deleted successfully",
	})
}

func DeleteAirdropFreeByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	err = module.DeleteAirdropFreeByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete AirdropFree by ID",
		})
	}

	invalidateAirdropCache()
	return c.JSON(fiber.Map{
		"message": "AirdropFree deleted successfully",
	})
}

func DeleteAirdropPaidByIDHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	err = module.DeleteAirdropPaidByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete AirdropPaid by ID",
		})
	}

	invalidateAirdropCache()
	return c.JSON(fiber.Map{
		"message": "AirdropPaid deleted successfully",
	})
}