package yahoo

type DinanceTsPoint struct {
	Time   int64   `json:"time"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

type YahooSymbolOCHL struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             *string  `json:"currency,omitempty"`
				Symbol               *string  `json:"symbol,omitempty"`
				ExchangeName         *string  `json:"exchangeName,omitempty"`
				FullExchangeName     *string  `json:"fullExchangeName,omitempty"`
				InstrumentType       *string  `json:"instrumentType,omitempty"`
				FirstTradeDate       *int     `json:"firstTradeDate,omitempty"`
				RegularMarketTime    *int     `json:"regularMarketTime,omitempty"`
				HasPrePostMarketData *bool    `json:"hasPrePostMarketData,omitempty"`
				Gmtoffset            *int     `json:"gmtoffset,omitempty"`
				Timezone             *string  `json:"timezone,omitempty"`
				ExchangeTimezoneName *string  `json:"exchangeTimezoneName,omitempty"`
				RegularMarketPrice   *float64 `json:"regularMarketPrice,omitempty"`
				FiftyTwoWeekHigh     *float64 `json:"fiftyTwoWeekHigh,omitempty"`
				FiftyTwoWeekLow      *float64 `json:"fiftyTwoWeekLow,omitempty"`
				RegularMarketDayHigh *float64 `json:"regularMarketDayHigh,omitempty"`
				RegularMarketDayLow  *float64 `json:"regularMarketDayLow,omitempty"`
				RegularMarketVolume  *int     `json:"regularMarketVolume,omitempty"`
				LongName             *string  `json:"longName,omitempty"`
				ShortName            *string  `json:"shortName,omitempty"`
				ChartPreviousClose   *float64 `json:"chartPreviousClose,omitempty"`
				PreviousClose        *float64 `json:"previousClose,omitempty"`
				Scale                *int     `json:"scale,omitempty"`
				PriceHint            *int     `json:"priceHint,omitempty"`
				CurrentTradingPeriod struct {
					Pre struct {
						Timezone  *string `json:"timezone,omitempty"`
						End       *int    `json:"end,omitempty"`
						Start     *int    `json:"start,omitempty"`
						Gmtoffset *int    `json:"gmtoffset,omitempty"`
					} `json:"pre"`
					Regular struct {
						Timezone  *string `json:"timezone,omitempty"`
						End       *int    `json:"end,omitempty"`
						Start     *int    `json:"start,omitempty"`
						Gmtoffset *int    `json:"gmtoffset,omitempty"`
					} `json:"regular"`
					Post struct {
						Timezone  *string `json:"timezone,omitempty"`
						End       *int    `json:"end,omitempty"`
						Start     *int    `json:"start,omitempty"`
						Gmtoffset *int    `json:"gmtoffset,omitempty"`
					} `json:"post"`
				} `json:"currentTradingPeriod"`
				TradingPeriods [][]struct {
					Timezone  *string `json:"timezone,omitempty"`
					End       *int    `json:"end,omitempty"`
					Start     *int    `json:"start,omitempty"`
					Gmtoffset *int    `json:"gmtoffset,omitempty"`
				} `json:"tradingPeriods,omitempty"`
				DataGranularity *string   `json:"dataGranularity,omitempty"`
				Range           *string   `json:"range,omitempty"`
				ValidRanges     []*string `json:"validRanges,omitempty"`
			} `json:"meta"`
			Timestamp  []*int `json:"timestamp,omitempty"`
			Indicators struct {
				Quote []struct {
					Close  []*float64 `json:"close,omitempty"`
					Open   []*float64 `json:"open,omitempty"`
					Volume []*int64   `json:"volume,omitempty"`
					Low    []*float64 `json:"low,omitempty"`
					High   []*float64 `json:"high,omitempty"`
				} `json:"quote,omitempty"`
			} `json:"indicators,omitempty"`
		} `json:"result,omitempty"`
		Error interface{} `json:"error,omitempty"`
	} `json:"chart,omitempty"`
}
