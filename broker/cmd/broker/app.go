package main

import (
	finnhub "broker/internal/api/finnhub"
	"broker/internal/api/yahoo"
	"broker/internal/config"
	cron "broker/internal/cron"
	market "broker/internal/market"
	"log"

	finnhub_api "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type App struct {
	Config         *config.Config
	Tickers        *config.Tickers
	MarketHolidays *market.MarketHolidays
	Cron           *cron.Cron
	FinnhubClient  *finnhub.FinnhubClient
	YahooClient    *yahoo.YahooClient
}

func NewApp() (*App, error) {
	appCfg, err := config.LoadConfig("configs/config.toml")
	if err != nil {
		return nil, err
	}
	log.Printf("Started %s in %s environment", appCfg.App.Name, appCfg.App.Env)

	tickers, err := config.LoadTickers("configs/tickers.json")
	if err != nil {
		return nil, err
	}

	holidays, err := market.LoadMarketHolidays("configs/holidays.json")
	if err != nil {
		return nil, err
	}

	c := cron.New()
	c.Start()

	finnhubCfg := finnhub_api.NewConfiguration()
	finnhubCfg.AddDefaultHeader("X-Finnhub-Token", appCfg.Finnhub.Token)
	apiClient := finnhub_api.NewAPIClient(finnhubCfg).DefaultApi

	return &App{
		Config:         appCfg,
		Tickers:        tickers,
		MarketHolidays: holidays,
		Cron:           c,
		FinnhubClient: &finnhub.FinnhubClient{
			Client:         apiClient,
			TickersConfig:  tickers,
			MarketHolidays: holidays,
		},
		YahooClient: &yahoo.YahooClient{
			TickersConfig: tickers,
		},
	}, nil
}

func (app *App) Run() {
	app.FinnhubClient.Run(app.Cron)
	app.YahooClient.Run(app.Cron)
}
