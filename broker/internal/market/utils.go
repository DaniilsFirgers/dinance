package market

import (
	"encoding/json"
	"fmt"
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
		return []HolidayRecord{} // Return empty slice for unsupported exchanges
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

func GetValidTradingPeriod(exchange Exchange, start time.Time, windowLength time.Duration, marketHolidays *MarketHolidays) (from, to time.Time, err error) {
	open, close, err := GetMarketHours(exchange)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if start.After(close) {
		return time.Time{}, time.Time{}, fmt.Errorf("period end %s is after market close %s", start, close)
	}

	isHoliday := marketHolidays.IsHoliday(exchange, start)
	if isHoliday {
		return time.Time{}, time.Time{}, fmt.Errorf("the date %s is a holiday for exchange %s", start.Format("2006-01-02"), exchange)
	}

	end := start.Add(-windowLength)

	if end.Before(open) {
		end = open
	}

	return end, start, nil
}
