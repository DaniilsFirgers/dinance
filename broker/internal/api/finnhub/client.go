package finnhub

import (
	"broker/internal/config"
	"broker/internal/market"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	cron "broker/internal/cron"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

const (
	REQUEST_DATE_FORMAT = "2006-01-02"
)

type FinnhubClient struct {
	Client         *finnhub.DefaultApiService
	TickersConfig  *config.Tickers
	MarketHolidays *market.MarketHolidays
}

func (f FinnhubClient) Run(cron *cron.Cron) error {
	cron.AddFunc("finnhub-company-news", "@every 15m", func() {
		if err := f.GetCompanyNews(); err != nil {
			log.Println("Error fetching company news:", err)
		}
	})
	return nil
}

func (f FinnhubClient) GetCompanyNews() error {
	var wg sync.WaitGroup
	from := time.Now().Format(REQUEST_DATE_FORMAT)
	to := time.Now().Format(REQUEST_DATE_FORMAT)

	tickers := append(f.TickersConfig.Tickers.EU, f.TickersConfig.Tickers.US...)

	for _, ticker := range tickers {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			res, _, err := f.Client.CompanyNews(context.Background()).Symbol(sym).From(from).To(to).Execute()
			if err != nil {
				fmt.Printf("Error fetching news for %s: %v\n", sym, err)
				return
			}
			fmt.Printf("Length of news for %s: %d\n", sym, len(res))
		}(ticker)
	}
	wg.Wait()
	return nil
}

func (f FinnhubClient) GetMarketNews() error {
	return nil
}
