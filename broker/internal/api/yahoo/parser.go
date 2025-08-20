package yahoo

import (
	"log"
	"time"
)

func checkPriceVolumeTrend(data YahooSymbolOCHL) {
	if data.Chart.Result == nil || len(data.Chart.Result) == 0 {
		log.Println("No chart data available")
		return
	}

	result := data.Chart.Result[0]
	if result.Timestamp == nil || len(result.Timestamp) == 0 {
		log.Printf("No timestamps available for symbol %s", *result.Meta.Symbol)
	}

	// Ensure we have at least one quote set
	if result.Indicators.Quote == nil || len(result.Indicators.Quote) == 0 {
		log.Printf("No quote data available for symbol %s", result.Meta.Symbol)
		return
	}

	quote := result.Indicators.Quote[0]
	if quote.Close == nil || quote.Volume == nil {
		log.Printf("No close or volume data available for symbol %s", *result.Meta.Symbol)
		return
	}

	cutoff := time.Now().Add(-6 * time.Hour).Unix()

	var points []DinanceTsPoint

	for i, ts := range result.Timestamp {
		if ts == nil {
			log.Printf("Nil timestamp at index %d for symbol %s", i, *result.Meta.Symbol)
			continue
		}

		if int64(*ts) < cutoff {
			continue
		}

		if quote.Close[i] == nil || quote.Volume[i] == nil {
			continue
		}

		points = append(points, DinanceTsPoint{
			Time:   int64(*ts),
			Close:  *quote.Close[i],
			Volume: *quote.Volume[i],
		})
	}

	if len(points) == 0 {
		log.Printf("No valid points found for symbol %s", *result.Meta.Symbol)
		return
	}

	// latest := points[len(points)-1]

}

func checkWindowTrend() {

}
