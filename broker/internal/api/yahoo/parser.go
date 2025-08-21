package yahoo

import (
	"log"
	"math"
	"time"
)

const MINIMUM_WINDOW = 5 * time.Minute
const DEFAULT_WINDOW_COUNT = 5

func checkPriceVolumeTrend(data YahooSymbolOCHL, cutoffTime time.Duration, windowsCount int) {
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

	cutoff := time.Now().Add(cutoffTime).Unix()

	var points []DinanceTsPoint

	for i, ts := range result.Timestamp {
		if ts == nil {
			log.Printf("Nil timestamp at index %d for symbol %s", i, *result.Meta.Symbol)
			continue
		}

		if *ts < cutoff {
			log.Printf("Skipping point for symbol %s at timestamp %d, before cutoff %d", *result.Meta.Symbol, *ts, cutoff)
			continue
		}

		if quote.Close[i] == nil || quote.Volume[i] == nil {
			log.Printf("Nil close or volume at index %d for symbol %s", i, *result.Meta.Symbol)
			continue
		}

		points = append(points, DinanceTsPoint{
			Time:   *ts,
			Close:  *quote.Close[i],
			Volume: *quote.Volume[i],
		})
	}

	if len(points) == 0 {
		log.Printf("No valid points found for symbol %s", *result.Meta.Symbol)
		return
	}
	windows := deriveWindows(cutoffTime, MINIMUM_WINDOW, windowsCount)
	computeWindowTrends(points, windows)
}

func computeWindowTrends(points []DinanceTsPoint, windows []time.Duration) {
	latest := points[len(points)-1]

	for _, window := range windows {
		ago := latest.Time - int64(window.Seconds())

		var past *DinanceTsPoint
		for _, p := range points {
			if p.Time >= ago {
				past = &p
				break
			}
		}

		if past == nil {
			log.Printf("No past point found for window %s", window)
			continue
		}

		priceChange := (latest.Close - past.Close) / past.Close * 100
		volumeChange := (latest.Volume - past.Volume) / past.Volume * 100

		log.Printf("Window %s: Price change: %.2f%%, Volume change: %.2f%%", window, priceChange, volumeChange)
	}
}

func deriveWindows(cuttOff, min time.Duration, count int) []time.Duration {
	if count <= 0 {
		return nil
	}

	if cuttOff < min {
		return nil
	}

	windows := make([]time.Duration, count)
	ratio := float64(cuttOff) / float64(min)
	step := math.Pow(ratio, 1/float64(count-1))

	cur := float64(min)
	for i := 0; i < count; i++ {
		windows[i] = time.Duration(cur).Round(time.Minute)
		cur *= step
	}
	return windows
}
