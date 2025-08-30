package yahoo

import (
	dinance_math "broker/internal/math"
	"fmt"
	"log"
	"math"
	"time"
)

const MINIMUM_DURATION = 3 * time.Minute
const DEFAULT_WINDOW_COUNT = 5

func checkPriceVolumeTrend(data YahooSymbolOCHL, start, end time.Time, maxDuration time.Duration) error {
	if data.Chart.Result == nil || len(data.Chart.Result) == 0 {
		return fmt.Errorf("no results found")
	}

	result := data.Chart.Result[0]
	if result.Timestamp == nil || len(result.Timestamp) == 0 {
		return fmt.Errorf("no timestamps found in result for symbol %s", *result.Meta.Symbol)
	}

	if result.Indicators.Quote == nil || len(result.Indicators.Quote) == 0 {
		return fmt.Errorf("no quote indicators found in result for symbol %s", *result.Meta.Symbol)
	}

	quote := result.Indicators.Quote[0]
	if quote.Close == nil || quote.Volume == nil {
		return fmt.Errorf("no close or volume data found in result for symbol %s", *result.Meta.Symbol)
	}

	dinance_math.ReverseSlicePtr(result.Timestamp)
	dinance_math.ReverseSlicePtr(result.Indicators.Quote[0].Volume)
	dinance_math.ReverseSlicePtr(result.Indicators.Quote[0].Close)

	startIndex, err := findStartingOchlTimeSeriesIndex(result)
	if err != nil {
		return err
	}

	var timeSeriesArray = result.Timestamp[startIndex:]
	var volumeArray = quote.Volume[startIndex:]
	var closeArray = quote.Close[startIndex:]

	// NOTE: will high certainty the last volume availble will be either 0 or null in Yahoo finance
	// therefore we cannot directly access it from the array
	var closePrice *float64 = closeArray[0]
	var closeVolume *int64 = volumeArray[0]
	var closeTime time.Time = time.Unix(*timeSeriesArray[0], 0).Truncate(time.Minute).UTC()

	windows := deriveWindowSteps(maxDuration, MINIMUM_DURATION, DEFAULT_WINDOW_COUNT, closeTime)

	var targets []DinanceTrendTsPoint
	cumulativeVolume := int64(0)

	// Timestamps are sorted in from oldest to newest, so we can iterate through them
	for i, ts := range timeSeriesArray {
		if ts == nil {
			continue
		}

		timestamp := time.Unix(*ts, 0).Truncate(time.Minute).UTC()

		if closeArray[i] == nil || volumeArray[i] == nil {
			continue
		}

		// NOTE if volume is zero, we skip the point as we will give untruthworthy changes (highly unlikely to have zero volume in a valid quote)
		if *volumeArray[i] == 0 {
			continue
		}

		cumulativeVolume += *volumeArray[i]
		if _, exists := windows[timestamp]; exists {
			targets = append(targets, DinanceTrendTsPoint{
				DinanceTsPoint: DinanceTsPoint{
					Time:   timestamp,
					Close:  *closeArray[i],
					Volume: *volumeArray[i],
				},
				Duration:         closeTime.Sub(timestamp),
				CumulativeVolume: cumulativeVolume,
			})

			delete(windows, timestamp)
		}
	}

	if len(targets) == 0 {
		return fmt.Errorf("no valid targets formed for symbol %s", *result.Meta.Symbol)
	}

	for _, target := range targets {
		priceChange := (*closePrice - target.Close) / target.Close * 100
		avgVolume := float64(target.CumulativeVolume) / closeTime.Sub(target.Time).Minutes()
		volumeRatio := float64(*closeVolume) / avgVolume
		fmt.Printf("Window: %s,Cum volume: %d, Avg window volume: %.2f, Duration: %f, Price change: %.2f%%, Volume ratio: %.2f\n", target.Time.Format(time.RFC3339), target.CumulativeVolume, avgVolume, target.Duration.Minutes(), priceChange, volumeRatio)
	}
	return nil
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

func findStartingOchlTimeSeriesIndex(timeSeries YahooSymbolOCHLResult) (int, error) {

	var validTsIndex int
	// iterate over timestamp to find first vaid volume and price
	for i, ts := range timeSeries.Timestamp {
		if ts == nil {
			continue
		}

		if timeSeries.Indicators.Quote[0].Volume[i] != nil && *timeSeries.Indicators.Quote[0].Volume[i] > 0 &&
			timeSeries.Indicators.Quote[0].Close[i] != nil && *timeSeries.Indicators.Quote[0].Close[i] > 0 {
			log.Printf("First valid timestamp: %s, unix: %d,  Volume: %d, Price: %.2f, Index: %d", time.Unix(*ts, 0).Format(time.RFC3339), *ts, *timeSeries.Indicators.Quote[0].Volume[i], *timeSeries.Indicators.Quote[0].Close[i], i)
			validTsIndex = i
			break
		}
	}

	if validTsIndex == 0 {
		return 0, fmt.Errorf("no valid timestamp found in time series for symbol %s", *timeSeries.Meta.Symbol)
	}

	return validTsIndex, nil
}
