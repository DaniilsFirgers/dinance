package config

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
)

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

func LoadConfig(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
