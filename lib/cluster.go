package lib

import (
	"math/rand"
	"time"
)

const (
	// minRate is the minimum rate at which we allow an action to be performed
	// across the whole cluster. The value is once a day: 1 / (1 * time.Day)
	minRate = 1.0 / 86400
)

// DurationMinusBuffer returns a duration, minus a buffer and jitter
// subtracted from the duration.  This function is used primarily for
// servicing Consul TTL Checks in advance of the TTL.
func DurationMinusBuffer(intv time.Duration, buffer time.Duration, jitter int64) time.Duration {
	d := intv - buffer
	if jitter == 0 {
		d -= RandomStagger(d)
	} else {
		d -= RandomStagger(time.Duration(int64(d) / jitter))
	}
	return d
}

// DurationMinusBufferDomain returns the domain of valid durations from a
// call to DurationMinusBuffer.  This function is used to check user
// specified input values to DurationMinusBuffer.
func DurationMinusBufferDomain(intv time.Duration, buffer time.Duration, jitter int64) (min time.Duration, max time.Duration) {
	max = intv - buffer
	if jitter == 0 {
		min = max
	} else {
		min = max - time.Duration(int64(max)/jitter)
	}
	return min, max
}

// RandomStagger returns an interval between 0 and the duration
func RandomStagger(intv time.Duration) time.Duration {
	if intv == 0 {
		return 0
	}
	return time.Duration(uint64(rand.Int63()) % uint64(intv))
}

// RandomStaggerWithRange returns an interval between min and the max duration
func RandomStaggerWithRange(min time.Duration, max time.Duration) time.Duration {
	return RandomStagger(max-min) + min
}

// RateScaledInterval is used to choose an interval to perform an action in
// order to target an aggregate number of actions per second across the whole
// cluster.
func RateScaledInterval(rate float64, min time.Duration, n int) time.Duration {
	if rate <= minRate {
		return min
	}
	interval := time.Duration(float64(time.Second) * float64(n) / rate)
	if interval < min {
		return min
	}

	return interval
}
