package yahoo

import (
	"time"
)

func getWindowMaxDuration(start, end time.Time, requestedMaxPeriod time.Duration) time.Duration {
	duration := end.Sub(start)
	if duration < requestedMaxPeriod {
		return duration
	}
	return requestedMaxPeriod
}
