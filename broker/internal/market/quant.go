package market

import (
	"fmt"
	"time"
)

const (
	US Exchange = "us"
	EU Exchange = "eu"
)

var marketHours = map[Exchange]MarketWorkingHours{
	US: {
		OpenHour:    13,
		OpenMinute:  30,
		CloseHour:   20,
		CloseMinute: 0,
		Timezone:    "America/New_York",
	},
	EU: {
		OpenHour:    10,
		OpenMinute:  0,
		CloseHour:   18,
		CloseMinute: 30,
		Timezone:    "Europe/Riga",
	},
}

func GetMarketHours(region Exchange) (start, end time.Time, err error) {
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		return time.Time{}, time.Time{}, fmt.Errorf("market is closed on weekends")
	}

	hours, exists := marketHours[region]
	if !exists {
		return time.Time{}, time.Time{}, fmt.Errorf("unknown market region: %s", region)
	}

	loc, err := time.LoadLocation(hours.Timezone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to load timezone %s: %v", hours.Timezone, err)
	}

	tNow := time.Now().In(loc)

	open := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), hours.OpenHour, hours.OpenMinute, 0, 0, loc)
	close := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), hours.CloseHour, hours.CloseMinute, 0, 0, loc)
	return open, close, nil

}
