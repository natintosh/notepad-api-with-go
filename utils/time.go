package utils

import (
	"time"
)

// GetTimeInMilliseconds :
var GetTimeInMilliseconds = func(t time.Time) int64 {
	return t.UnixNano() / int64(time.Second) * int64(time.Microsecond)
}

// GetTimeInGo :
var GetTimeInGo = func(i int64) time.Time {
	t := i * int64(time.Millisecond) / int64(time.Second)
	return time.Unix(t, 0)
}
