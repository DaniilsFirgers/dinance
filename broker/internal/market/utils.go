package market

import (
	"encoding/json"
	"os"
	"time"
)

const (
	HOLIDAY_DATE_FORMAT = "2006-01-02"
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

func (m *MarketHolidays) holidaysForExchange(exchange Exchange) []HolidayRecord {
	switch exchange {
	case US:
		return m.US
	case EU:
		return m.EU
	default:
		return nil
	}
}

func (m *MarketHolidays) IsHoliday(exchange Exchange, date time.Time) bool {
	dateStr := date.Format(HOLIDAY_DATE_FORMAT)
	holidays := m.holidaysForExchange(exchange)
	for _, holiday := range holidays {
		if holiday.Date == dateStr {
			return true
		}
	}
	return false
}
