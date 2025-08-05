package configs

import (
	"encoding/json"
	"os"
)

type Tickers struct {
	Tickers []string `json:"tickers"` // List of stock tickers.
}

func LoadTickers(path string) (*Tickers, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tickers Tickers
	err = json.Unmarshal(file, &tickers)
	if err != nil {
		return nil, err
	}

	return &tickers, nil
}
