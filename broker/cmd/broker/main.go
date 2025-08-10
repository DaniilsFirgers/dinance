package main

import (
	finnhub "broker/internal/api/finnhub"
	"broker/internal/config"
	"log"

	finnhub_api "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/robfig/cron/v3"
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
	cron := cron.New(cron.WithSeconds())
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
