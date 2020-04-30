package time_test

import (
  "github.com/stretchr/testify/assert"
	testtime "github.com/thecodedproject/gotest/time"
	"testing"
	"time"
)

func TestSetTimeNowFuncAndReset(t *testing.T) {

	someTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	reset := testtime.SetTimeNowFuncForTesting(t, func() time.Time {
		return someTime
	})
	assert.Equal(t, someTime, testtime.Now())
	assert.Equal(t, someTime, testtime.Now())

	reset()

	now := testtime.Now()
	diff := time.Now().Sub(now)
	assert.True(t, diff < 5*time.Millisecond)
}

func TestSetTimeNowValueAndReset(t *testing.T) {

	someTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	reset := testtime.SetTimeNowForTesting(t, someTime)

	assert.Equal(t, someTime, testtime.Now())
	assert.Equal(t, someTime, testtime.Now())

	reset()

	now := testtime.Now()
	diff := time.Now().Sub(now)
	assert.True(t, diff < 5*time.Millisecond)
}

func TestTimeNowWithoutMocking(t *testing.T) {

	now := testtime.Now()
	diff := time.Now().Sub(now)
	assert.True(t, diff < 5*time.Millisecond)
}
