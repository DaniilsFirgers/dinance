package yahoo

import (
	"broker/internal/config"
	cron "broker/internal/cron"
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

	res, err := httpclient.Get("https://query1.finance.yahoo.com/v8/finance/chart/AAPL?interval=1m&range=1d", httpclient.GetHeaders("https://finance.yahoo.com/"))
	if err != nil {
		log.Println("Error fetching quote:", err)
		return "", err
	}

	var news YahooSymbolOCHL
	if err := json.Unmarshal(res, &news); err != nil {
		log.Println("Error unmarshaling response:", err)
		return "", err
	}

	return "Quote for " + symbol, nil
}
