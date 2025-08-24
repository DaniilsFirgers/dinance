package yahoo

import (
	"log"
	"math"
	"time"
)

const MINIMUM_WINDOW = 5 * time.Minute
const DEFAULT_WINDOW_COUNT = 5

func checkPriceVolumeTrend(data YahooSymbolOCHL, maxWindow time.Duration, windowsCount int) {
	if data.Chart.Result == nil || len(data.Chart.Result) == 0 {
		return
	}

	result := data.Chart.Result[0]
	if result.Timestamp == nil || len(result.Timestamp) == 0 {
		return
	}

	if result.Indicators.Quote == nil || len(result.Indicators.Quote) == 0 {
		return
	}

	quote := result.Indicators.Quote[0]
	if quote.Close == nil || quote.Volume == nil {
		return
	}

	cutoff := time.Now().Truncate(time.Minute).Add(-maxWindow)

	var points []DinanceTsPoint
	// Timestamps are sorted in from oldest to newest, so we can iterate through them
	for i, ts := range result.Timestamp {
		if ts == nil {
			continue
		}

		timestamp := time.Unix(*ts, 0).Truncate(time.Minute).UTC()
		if timestamp.Before(cutoff) {
			continue
		}

		if quote.Close[i] == nil || quote.Volume[i] == nil {
			continue
		}

		// NOTE if volume is zero, we skip the point as we will give untruthworthy changes (highly unlikely to have zero volume in a valid quote)
		if *quote.Volume[i] == 0 {
			continue
		}

		points = append(points, DinanceTsPoint{
			Time:   timestamp,
			Close:  *quote.Close[i],
			Volume: *quote.Volume[i],
		})
	}

	if len(points) == 0 {
		log.Printf("No valid points found for symbol %s", *result.Meta.Symbol)
		return
	}
	computeWindowTrends(points, maxWindow)
}

func computeWindowTrends(points []DinanceTsPoint, maxWindow time.Duration) {
	latest := points[len(points)-1]

	windows := deriveWindowSteps(maxWindow, MINIMUM_WINDOW, DEFAULT_WINDOW_COUNT, latest.Time)
	for w, _ := range windows {
		log.Printf("Window: %s", w.Format(time.RFC3339))
	}
	var targets []DinanceTrendTsPoint
	cumulativeVolume := int64(0)

	for _, p := range points {
		cumulativeVolume += p.Volume

		if _, exists := windows[p.Time]; exists {
			targets = append(targets, DinanceTrendTsPoint{
				DinanceTsPoint: DinanceTsPoint{
					Time:   p.Time,
					Close:  p.Close,
					Volume: p.Volume,
				},
				CumulativeVolume: cumulativeVolume,
			})

			delete(windows, p.Time)
			break
		}
	}

	if len(targets) == 0 {
		return
	}
	for _, target := range targets {
		priceChange := (latest.Close - target.Close) / target.Close * 100
		avgVolume := float64(target.CumulativeVolume) / latest.Time.Sub(target.Time).Minutes()
		volumeRatio := float64(latest.Volume) / avgVolume
		log.Printf("Window %s: Price change: %.2f%%, Volume ratio: %.2f", target.Time.Format(time.RFC3339), priceChange, volumeRatio)
	}
}

func deriveWindowSteps(maxDuration, minDuration time.Duration, count int, targetTime time.Time) map[time.Time]struct{} {
	if count <= 0 {
		return nil
	}

	if maxDuration < minDuration {
		return nil
	}

	windowSteps := make(map[time.Time]struct{}, count)
	ratio := float64(maxDuration) / float64(minDuration)
	step := math.Pow(ratio, 1/float64(count-1))

	currDuration := minDuration
	for i := 0; i < count; i++ {
		stepTime := targetTime.Add(-currDuration)
		if _, exists := windowSteps[stepTime]; !exists {
			windowSteps[stepTime] = struct{}{}
		}

		currDuration = time.Duration(float64(currDuration) * step)
	}
	return windowSteps
}
