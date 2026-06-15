package models

type CryptoResponse struct {
	Data     map[string]CryptoData `json:"data"`
	Metadata Metadata              `json:"metadata"`
}

type CryptoData struct {
	ID                int              `json:"id"`
	Name              string           `json:"name"`
	Symbol            string           `json:"symbol"`
	WebsiteSlug       string           `json:"website_slug"`
	Rank              int              `json:"rank"`
	CirculatingSupply float64          `json:"circulating_supply"`
	TotalSupply       float64          `json:"total_supply"`
	MaxSupply         float64          `json:"max_supply"`
	Quotes            map[string]Quote `json:"quotes"`
	LastUpdated       int64            `json:"last_updated"`
}

type Quote struct {
	Price               float64 `json:"price"`
	Volume24h           float64 `json:"volume_24h"`
	MarketCap           float64 `json:"market_cap"`
	PercentageChange1h  float64 `json:"percentage_change_1h"`
	PercentageChange24h float64 `json:"percentage_change_24h"`
	PercentageChange7d  float64 `json:"percentage_change_7d"`
	PercentChange1h     float64 `json:"percent_change_1h"`
	PercentChange24h    float64 `json:"percent_change_24h"`
	PercentChange7d     float64 `json:"percent_change_7d"`
}

type Metadata struct {
	Timestamp           int64       `json:"timestamp"`
	NumCryptocurrencies int         `json:"num_cryptocurrencies"`
	Error               interface{} `json:"error"`
}