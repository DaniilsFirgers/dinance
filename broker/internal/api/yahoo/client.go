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

	job := func() {
		periodStart := time.Now().UTC().Truncate(time.Minute)
		start, end, err := getRequestPeriods(periodStart, 6*time.Hour)
		if err != nil {
			log.Printf("Error getting request periods: %v", err)
			return
		}

		maxDuration := getWindowMaxDuration(start, end, 6*time.Hour)
		if err := y.GetQuotesData(start, end, maxDuration); err != nil {
			log.Println("Error fetching quotes data:", err)
		}
	}

	// First cron: 10:00–17:59 EET, every 2 minutes
	cron.AddFunc("yahoo-quote-part-one", "0-59/2 10-17 * * 1-5", job)
	// Second cron: 18:00–18:30 EET, every 2 minutes
	cron.AddFunc("yahoo-quote-part-two", "0-30/2 18 * * 1-5", job)
}

func (y YahooClient) GetQuotesData(start, end time.Time, maxDuration time.Duration) error {
	var wg sync.WaitGroup

	for _, symbol := range y.TickersConfig.Tickers {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			data, err := y.GetQuoteData(sym, start, end)
			if err != nil {
				log.Printf("Error fetching quote for %s: %v\n", sym, err)
				return
			}
			if err := checkPriceVolumeTrend(data, start, end, maxDuration); err != nil {
				log.Printf("Error checking price volume trend for %s: %v\n", sym, err)
			}
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
