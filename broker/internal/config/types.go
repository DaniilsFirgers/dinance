package config

type Tickers struct {
	Tickers []string `json:"tickers"` // List of stock tickers.
}

type Config struct {
	App struct {
		Name string
		Env  string
	}
	Finnhub struct {
		Token string
	}
}
