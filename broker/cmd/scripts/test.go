package main

import (
	"broker/internal/api/yahoo"
	"broker/internal/config"
	"fmt"
	"log"
	"time"
)

func main() {
	tickers, err := config.LoadTickers("configs/tickers.json")
	if err != nil {
		log.Fatal("Error loading tickers:", err)
	}

	yahooClient := yahoo.YahooClient{TickersConfig: tickers}

	// Last friday's date
	start := time.Date(2025, time.August, 29, 15, 0, 0, 0, time.UTC)
	end := time.Date(2025, time.August, 29, 19, 0, 0, 0, time.UTC)
	fmt.Printf("Fetching quotes data from %s to %s\n", start, end)

	if err := yahooClient.GetQuotesData(start, end, (3 * time.Hour)); err != nil {
		log.Println("Error fetching quotes data:", err)
	}
}
