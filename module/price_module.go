package module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"

	"github.com/nekowawolf/airdropv2/models"
)

func GetPrice(url string) (*models.CryptoData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response models.CryptoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	for _, v := range response.Data {
		return &v, nil
	}

	return nil, fmt.Errorf("data not found")
}