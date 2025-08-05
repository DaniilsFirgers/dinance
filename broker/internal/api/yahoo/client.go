package yahoo

import (
	httpclient "broker/internal/http-client"
	"encoding/json"
	"fmt"
	"log"
)

func GetQuote(symbol string) (string, error) {
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0 Safari/537.36",
		"Accept":          "application/json, text/javascript, */*; q=0.01",
		"Accept-Language": "en-US,en;q=0.9",
		"Referer":         "https://finance.yahoo.com/",
		"Connection":      "keep-alive",
	}

	res, err := httpclient.Get("https://query1.finance.yahoo.com/v8/finance/chart/AAPL?interval=1m&range=1d", headers)
	if err != nil {
		log.Println("Error fetching quote:", err)
		return "", err
	}

	var news YahooSymbolOCHL
	if err := json.Unmarshal(res, &news); err != nil {
		log.Println("Error unmarshaling response:", err)
		return "", err
	}
	jsonBytes, err := json.MarshalIndent(news, "", "  ")
	if err != nil {
		log.Println("Error marshaling:", err)
	} else {
		fmt.Println(string(jsonBytes))
	}
	return "Quote for " + symbol, nil
}
