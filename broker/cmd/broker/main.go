package main

import (
	finnhub "broker/internal/api/finnhub"
	"broker/internal/config"
	"log"

	cron "broker/internal/cron"

	finnhub_api "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func main() {

	app_cfg, err := config.LoadConfig("configs/config.toml")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	log.Printf("Started %s in %s environment", app_cfg.App.Name, app_cfg.App.Env)

	// Load tickers I am interested in
	tickers, err := config.LoadTickers("configs/tickers.json")
	if err != nil {
		log.Fatal("Error loading tickers:", err)
	}

	// Initialize the cron scheduler and start it
	cron := cron.New()
	cron.Start()

	finnhub_cfg := finnhub_api.NewConfiguration()
	finnhub_cfg.AddDefaultHeader("X-Finnhub-Token", app_cfg.Finnhub.Token)
	apiClient := finnhub_api.NewAPIClient(finnhub_cfg).DefaultApi

	finnhub_client := finnhub.FinnhubClient{Client: apiClient, TickersConfig: tickers}

	finnhub_client.Run(cron)
	// yahoo.GetQuote("GOOGL")
	// Keep the program running
	select {}
}
