package config

type Tickers struct {
	Tickers struct {
		US []string `json:"us"`
		EU []string `json:"eu"`
	} `json:"tickers"`
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
