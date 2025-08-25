package yahoo

import (
	"broker/internal/config"
	cron "broker/internal/cron"
	"fmt"
	"sync"
	"time"

	httpclient "broker/internal/http-client"
	"encoding/json"
	"log"
)

type YahooClient struct {
	TickersConfig *config.Tickers
}

func (y YahooClient) Run(cron *cron.Cron) {
	if err := y.GetQuotesData(time.Hour * 6); err != nil {
		log.Println("Error fetching quotes data:", err)
	}
	// cron.AddFunc("yahoo-quote", "@every 10s", func() {
	// 	if err := y.GetQuotesData(time.Hour * 3); err != nil {
	// 		log.Println("Error fetching quotes data:", err)
	// 	}
	// })
}

func (y YahooClient) GetQuotesData(maxDuration time.Duration) error {
	var wg sync.WaitGroup

	periodStart := time.Now().UTC().Truncate(time.Minute)
	start, end, err := getRequestPeriods(periodStart, maxDuration)
	if err != nil {
		return fmt.Errorf("failed to get request periods: %w", err)
	}
	windowMaxDuration := getWindowMaxDuration(start, end, maxDuration)

	for _, symbol := range y.TickersConfig.Tickers {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			data, err := y.GetQuoteData(sym, start, end)
			if err != nil {
				log.Printf("Error fetching quote for %s: %v\n", sym, err)
				return
			}
			checkPriceVolumeTrend(data, windowMaxDuration)
		}(symbol)
	}
	wg.Wait()
	return nil
}

func (y YahooClient) GetQuoteData(symbol string, startTime, endTime time.Time) (YahooSymbolOCHL, error) {
	headers := httpclient.GetHeaders("https://finance.yahoo.com/")
	fmt.Printf("Fetching data for symbol: %s from %s to %s\n", symbol, startTime.Unix(), endTime.Unix())
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1m&range=1d&period1=%d&period2=%d", symbol, startTime.Unix(), endTime.Unix())

	res, err := httpclient.Get(url, headers)
	if err != nil {
		return YahooSymbolOCHL{}, err
	}

	var data YahooSymbolOCHL
	if err := json.Unmarshal(res, &data); err != nil {
		return YahooSymbolOCHL{}, err
	}

	return data, nil
}
