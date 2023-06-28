package time

import (
	"testing"
	gotime "time"
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
