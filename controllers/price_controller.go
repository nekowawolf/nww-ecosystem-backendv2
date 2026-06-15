package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nekowawolf/airdropv2/models"
	"github.com/nekowawolf/airdropv2/module"
	"github.com/nekowawolf/airdropv2/utils"
)

func PriceHandler(c *fiber.Ctx) error {
	coins := map[string]string{
		"btc":   "https://api.alternative.me/v2/ticker/bitcoin/",
		"eth":   "https://api.alternative.me/v2/ticker/ethereum/",
		"sol":   "https://api.alternative.me/v2/ticker/solana/",
		"bnb":   "https://api.alternative.me/v2/ticker/binancecoin/",
		"matic": "https://api.alternative.me/v2/ticker/matic-network/",
		"xrp":   "https://api.alternative.me/v2/ticker/ripple/",
	}

	results := make(map[string]interface{})

	for key, url := range coins {
		currentURL := url
		data, err := utils.GetOrSetCache("price:"+key, 5*time.Minute, func() (*models.CryptoData, error) {
			return module.GetPrice(currentURL)
		})
		
		if err != nil {
			log.Println("Error fetching price for", key, ":", err)
			results[key] = "Error"
		} else {
			results[key] = data
		}
	}

	return c.JSON(results)
}