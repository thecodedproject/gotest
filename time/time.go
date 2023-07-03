package time

import (
	"testing"
	gotime "time"
)

const (
	Nanosecond  = gotime.Nanosecond
	Microsecond = gotime.Microsecond
	Millisecond = gotime.Millisecond
	Second = gotime.Second
	Minute = gotime.Minute
	Hour = gotime.Hour
)

var nowFunc = gotime.Now

func Now() gotime.Time {

	return nowFunc()
}

func SetTimeNowFuncForTesting(t *testing.T, now func() gotime.Time) {

	oldNowFunc := nowFunc
	nowFunc = now
	t.Cleanup(func() {
		nowFunc = oldNowFunc
	})
}

func SetTimeNowForTesting(t *testing.T) (now gotime.Time) {

	now = gotime.Now()

	SetTimeNowFuncForTesting(t, func() gotime.Time {
		return now
	})

	return now
}
