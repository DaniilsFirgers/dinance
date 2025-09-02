package config

type Tickers struct {
	Tickers struct {
		US []string `json:"us"`
		EU []string `json:"eu"`
	} `json:"tickers"`
}

type Config struct {
	App struct {
		Name string `json:"name"`
		Env  string `json:"env"`
	} `json:"app"`
	Finnhub struct {
		Token string `json:"token"`
	} `json:"finnhub"`
	Settings ConfigSettings `json:"settings"`
}

type ConfigSettings struct {
	Trend struct {
		VolumeRatio int `json:"volume_ratio"`
	} `json:"trend"`
}
