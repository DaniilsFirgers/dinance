package finnhub

import (
	"broker/configs"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

const (
	REQUEST_DATE_FORMAT = "2006-01-02"
)

type FinnhubClient struct {
	Client *finnhub.DefaultApiService
}

func (f FinnhubClient) GetCompanyNews(config configs.Tickers) error {
	var wg sync.WaitGroup
	from := time.Now().Format(REQUEST_DATE_FORMAT)
	to := time.Now().Format(REQUEST_DATE_FORMAT)

	for _, symbol := range config.Tickers {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			res, _, err := f.Client.CompanyNews(context.Background()).Symbol(sym).From(from).To(to).Execute()
			if err != nil {
				fmt.Printf("Error fetching news for %s: %v\n", sym, err)
				return
			}
			jsonBytes, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				log.Println("Error marshaling:", err)
			} else {
				fmt.Println(string(jsonBytes))
			}
		}(symbol)
	}
	wg.Wait()
	return nil
}

func (f FinnhubClient) GetMarketNews() error {
	// from := time.Now().Format(REQUEST_DATE_FORMAT)
	// to := time.Now().Format(REQUEST_DATE_FORMAT)

	return nil
}
