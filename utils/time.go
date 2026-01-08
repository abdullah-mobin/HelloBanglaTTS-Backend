package utils

import (
	"fmt"
	"time"
)

// NowUnixString returns the current Unix timestamp as a string
func NowUnixString() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func Now() time.Time {
	return time.Now()
}
