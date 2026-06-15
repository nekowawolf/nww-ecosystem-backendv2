package controllers

import (
	"log"
	"sync"
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
	var wg sync.WaitGroup
	var mu sync.Mutex

	for key, url := range coins {
		wg.Add(1)
		go func(k string, u string) {
			defer wg.Done()
			
			data, err := utils.GetOrSetCache("price:"+k, 5*time.Minute, func() (*models.CryptoData, error) {
				return module.GetPrice(u)
			})
			
			mu.Lock()
			defer mu.Unlock()
			
			if err != nil {
				log.Println("Error fetching price for", k, ":", err)
				results[k] = "Error"
			} else {
				results[k] = data
			}
		}(key, url)
	}

	wg.Wait()

	return c.JSON(results)
}