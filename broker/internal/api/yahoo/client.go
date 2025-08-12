package yahoo

import (
	"broker/internal/config"
	cron "broker/internal/cron"
	"fmt"
	"sync"

	httpclient "broker/internal/http-client"
	"encoding/json"
	"log"
)

type YahooClient struct {
	TickersConfig *config.Tickers
}

func (y YahooClient) Run(cron *cron.Cron) {
	cron.AddFunc("yahoo-quote", "@every 10s", func() {
		if err := y.GetQuotesData(); err != nil {
			log.Println("Error fetching quotes data:", err)
		}
	})
}

func (y YahooClient) GetQuotesData() error {
	var wg sync.WaitGroup

	for _, symbol := range y.TickersConfig.Tickers {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			quote, err := y.getQuote(sym)
			if err != nil {
				log.Printf("Error fetching quote for %s: %v\n", sym, err)
				return
			}
			log.Printf("Quote for %s: %s\n", sym, quote)
		}(symbol)
	}
	wg.Wait()
	return nil
}

func (y YahooClient) getQuote(symbol string) (string, error) {
	headers := httpclient.GetHeaders("https://finance.yahoo.com/")
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1m&range=1d", symbol)

	res, err := httpclient.Get(url, headers)
	if err != nil {
		log.Println("Error fetching quote:", err)
		return "", err
	}

	var data YahooSymbolOCHL
	if err := json.Unmarshal(res, &data); err != nil {
		log.Println("Error unmarshaling response:", err)
		return "", err
	}

	return "Quote for " + symbol, nil
}
