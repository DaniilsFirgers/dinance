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

func getRequestPeriods(periodEnd time.Time, windowLength time.Duration) (from, to time.Time, err error) {
	open, close, error := market.GetMarketHours()
	if error != nil {
		return time.Time{}, time.Time{}, error
	}

	if periodEnd.After(close) {
		return time.Time{}, time.Time{}, fmt.Errorf("period end %s is after market close %s", periodEnd, close)
	}

	start := periodEnd.Add(-windowLength)

	if start.Before(open) {
		start = open
	}

	return start, periodEnd, nil
}
