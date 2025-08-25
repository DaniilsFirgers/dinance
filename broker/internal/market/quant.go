package market

import (
	"fmt"
	"time"
)

func GetMarketHours() (start, end time.Time, err error) {
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		return time.Time{}, time.Time{}, fmt.Errorf("market is closed on weekends")
	}
	open := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 13, 30, 0, 0, time.UTC)
	close := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 20, 0, 0, 0, time.UTC)
	return open, close, nil

}
