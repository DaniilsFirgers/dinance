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
	cron.AddFunc("yahoo-quote", "@every 10s", func() {
		if err := y.GetQuotesData(time.Hour * 3); err != nil {
			log.Println("Error fetching quotes data:", err)
		}
	})
}

func (y YahooClient) GetQuotesData(cutOffTime time.Duration) error {
	var wg sync.WaitGroup

	for _, symbol := range y.TickersConfig.Tickers {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			data, err := y.getQuoteData(sym)
			if err != nil {
				log.Printf("Error fetching quote for %s: %v\n", sym, err)
				return
			}
			checkPriceVolumeTrend(data, cutOffTime, DEFAULT_WINDOW_COUNT)
		}(symbol)
	}
	wg.Wait()
	return nil
}

func (y YahooClient) getQuoteData(symbol string) (YahooSymbolOCHL, error) {
	headers := httpclient.GetHeaders("https://finance.yahoo.com/")
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1m&range=1d", symbol)

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
