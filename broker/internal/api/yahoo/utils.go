package yahoo

import (
	"broker/internal/market"
	"fmt"
	"time"
)

func getWindowMaxDuration(start, end time.Time, requestedMaxPeriod time.Duration) time.Duration {
	duration := end.Sub(start)
	if duration < requestedMaxPeriod {
		return duration
	}
	return requestedMaxPeriod
}

func getRequestPeriods(exchange market.Exchange, start time.Time, windowLength time.Duration, marketHolidays *market.MarketHolidays) (from, to time.Time, err error) {
	open, close, err := market.GetMarketHours(exchange)
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
