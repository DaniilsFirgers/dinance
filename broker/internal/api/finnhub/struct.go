package finnhub

type CompanyNews struct {
	Category string `json:"category"` // News category.
	Datetime int64  `json:"datetime"` // Unix timestamp of when the news was published.
	Headline string `json:"headline"` // News headline.
	ID       int64  `json:"id"`       // Unique identifier for the news article. This can be used with the `minId` parameter to fetch only the latest news articles.
	Image    string `json:"image"`    // URL of the thumbnail image associated with the news article.
	Related  string `json:"related"`  // Stocks and companies related to the news article.
	Source   string `json:"source"`   // Source of the news article.
	Summary  string `json:"summary"`  // Summary of the news article.
	URL      string `json:"url"`      // URL of the original article.
}

type MarketNewsCategory string

const (
	CategoryGeneral MarketNewsCategory = "general"
	CategoryForex   MarketNewsCategory = "forex"
	CategoryCrypto  MarketNewsCategory = "crypto"
	CategoryMerger  MarketNewsCategory = "merger"
)
