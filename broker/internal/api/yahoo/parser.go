package yahoo

import (
	"fmt"
	"log"
	"math"
	"time"
)

const MINIMUM_DURATION = 3 * time.Minute
const DEFAULT_WINDOW_COUNT = 5

func checkPriceVolumeTrend(data YahooSymbolOCHL, maxDuration time.Duration) {
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

	lastTimeSec := result.Timestamp[len(result.Timestamp)-1]
	if lastTimeSec == nil {
		log.Printf("No valid timestamp found for symbol %s", *result.Meta.Symbol)
		return
	}

	lastTime := time.Unix(*lastTimeSec, 0).Truncate(time.Minute).UTC()
	lastPrice := quote.Close[len(quote.Close)-1]
	lastVolume := quote.Volume[len(quote.Volume)-1]

	windows := deriveWindowSteps(maxDuration, MINIMUM_DURATION, DEFAULT_WINDOW_COUNT, lastTime)

	var targets []DinanceTrendTsPoint
	cumulativeVolume := int64(0)

	// Timestamps are sorted in from oldest to newest, so we can iterate through them
	for i, ts := range result.Timestamp {
		if ts == nil {
			continue
		}

		timestamp := time.Unix(*ts, 0).Truncate(time.Minute).UTC()

		if quote.Close[i] == nil || quote.Volume[i] == nil {
			continue
		}

		// NOTE if volume is zero, we skip the point as we will give untruthworthy changes (highly unlikely to have zero volume in a valid quote)
		if *quote.Volume[i] == 0 {
			continue
		}

		cumulativeVolume += *quote.Volume[i]

		if _, exists := windows[timestamp]; exists {
			targets = append(targets, DinanceTrendTsPoint{
				DinanceTsPoint: DinanceTsPoint{
					Time:   timestamp,
					Close:  *quote.Close[i],
					Volume: *quote.Volume[i],
				},
				CumulativeVolume: cumulativeVolume,
			})

			delete(windows, timestamp)
		}
	}

	if len(targets) == 0 {
		log.Printf("No valid targets formed for symbol %s", *result.Meta.Symbol)
		return
	}

	for _, target := range targets {
		priceChange := (*lastPrice - target.Close) / target.Close * 100
		avgVolume := float64(target.CumulativeVolume) / lastTime.Sub(target.Time).Minutes()
		fmt.Printf("VOlume: %d, Avg Volume: %.2f, Last Volume: %d\n", target.CumulativeVolume, avgVolume, *lastVolume)
		volumeRatio := float64(*lastVolume) / avgVolume
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

		currDuration = time.Duration(float64(currDuration) * step).Truncate(time.Minute)
	}
	return windowSteps
}
