package market

type Exchange string

type MarketWorkingHours struct {
	OpenHour    int
	CloseHour   int
	OpenMinute  int
	CloseMinute int
	Timezone    string
}
