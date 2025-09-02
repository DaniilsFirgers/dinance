package market

import (
	"encoding/json"
	"os"
)

func LoadMarketHolidays(path string) (*MarketHolidays, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var holidays MarketHolidays
	err = json.Unmarshal(file, &holidays)
	if err != nil {
		return nil, err
	}

	return &holidays, nil
}
