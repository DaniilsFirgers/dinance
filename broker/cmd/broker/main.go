package main

import (
	"broker/configs"
	finnhub "broker/internal/api/finnhub"
	"log"

	finnhub_api "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func main() {

	tickers, err := configs.LoadTickers("configs/tickers.json")
	if err != nil {
		log.Fatal("Error loading tickers:", err)
	}

	cfg := finnhub_api.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", "")
	apiClient := finnhub_api.NewAPIClient(cfg).DefaultApi

	finnhub_client := finnhub.FinnhubClient{Client: apiClient}

	finnhub_client.GetCompanyNews(*tickers)
	// yahoo.GetQuote("GOOGL")
	// Keep the program running
	select {}
}
