package market

type Exchange string

type MarketWorkingHours struct {
	OpenHour    int
	CloseHour   int
	OpenMinute  int
	CloseMinute int
	Timezone    string
}

type MarketHolidays struct {
	US []HolidayRecord `json:"us"`
	EU []HolidayRecord `json:"eu"`
}

type HolidayRecord struct {
	Date        string `json:"date"`
	EventName   string `json:"event_name"`
	TradingHour string `json:"trading_hour"`
}
