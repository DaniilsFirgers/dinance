package yahoo

import (
	"broker/internal/config"
	cron "broker/internal/cron"
	"fmt"
	"sync"
	"time"

	httpclient "broker/internal/http-client"
	market "broker/internal/market"
	"encoding/json"
	"log"
)

const (
	DEFAULT_WINDOW_LENGTH = 3 * time.Hour
)

type YahooClient struct {
	TickersConfig  *config.Tickers
	MarketHolidays *market.MarketHolidays
}

func (y YahooClient) Run(cron *cron.Cron) {
	job := func(exchange market.Exchange, requestPeriod time.Duration) {
		periodStart := time.Now().UTC().Truncate(time.Minute)
		start, end, err := getRequestPeriods(exchange, periodStart, requestPeriod, y.MarketHolidays)
		if err != nil {
			log.Printf("Error getting request periods: %v", err)
			return
		}

		maxDuration := getWindowMaxDuration(start, end, requestPeriod)
		if err := y.GetQuotesData(exchange, start, end, maxDuration); err != nil {
			log.Println("Error fetching quotes data:", err)
		}
	}

	// First cron: 10:00–17:59 EET, every 2 minutes for EU region
	cron.AddFunc("yahoo-quote-part-one", "0-59/2 10-17 * * 1-5", func() {
		job(market.EU, DEFAULT_WINDOW_LENGTH)
	})
	// Second cron: 18:00–18:30 EET, every 2 minutes for EU region
	cron.AddFunc("yahoo-quote-part-two", "0-30/2 18 * * 1-5", func() {
		job(market.EU, DEFAULT_WINDOW_LENGTH)
	})

	// Third cron: 16:30–16:59 EET, every 2 minutes for US region
	cron.AddFunc("yahoo-quote-part-three", "30-59/2 16 * * 1-5", func() {
		job(market.US, time.Minute*30) // Shorter window for the first half hour
	})

	// Fourth cron: 17:00-22:59 EET, every 2 minutes for US region
	cron.AddFunc("yahoo-quote-part-four", "*/2 17-22 * * 1-5", func() {
		job(market.US, DEFAULT_WINDOW_LENGTH)
	})
}

func (y YahooClient) GetQuotesData(exchange market.Exchange, start, end time.Time, maxDuration time.Duration) error {

	parser := func(tickers []string) {
		var wg sync.WaitGroup

		for _, ticker := range tickers {
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
			}(ticker)
		}
		wg.Wait()
	}

	switch exchange {
	case market.US:
		parser(y.TickersConfig.Tickers.US)
	case market.EU:
		parser(y.TickersConfig.Tickers.EU)
	default:
		return fmt.Errorf("unsupported region: %s", exchange)
	}

	return nil
}

func (y YahooClient) GetQuoteData(symbol string, startTime, endTime time.Time) (YahooSymbolOCHL, error) {
	headers := httpclient.GetHeaders("https://finance.yahoo.com/")
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
