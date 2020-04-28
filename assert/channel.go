package assert

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

const assertionTimeout = 100*time.Millisecond

func ChannelReceivesOnce(
	t *testing.T,
	ch <-chan interface{},
	expected interface{},
) func() bool {

	var wg sync.WaitGroup
	wg.Add(1)

	result := true

	go func() {
		select {
			case v := <-ch:
				result = result && assert.Equal(t, v, expected)
			case <-time.After(assertionTimeout):
				assert.Fail(t, "Channel recieved nothing", "Expected:", expected)
				result = false
		}
		wg.Done()
	}()

	return func() bool {
		wg.Wait()
		return result
	}
}

func ChannelReceives(
	t *testing.T,
	ch <-chan interface{},
	expected []interface{},
) func() bool {

	var wg sync.WaitGroup
	wg.Add(1)

	var result bool

	go func() {

		recieved := make([]interface{}, 0)
		for {
			select {
				case v := <-ch:
					recieved = append(recieved, v)
				case <-time.After(200*time.Millisecond):
					result = assert.Equal(t, recieved, expected)
					wg.Done()
					return
			}
		}
	}()

	return func() bool {
		wg.Wait()
		return result
	}
}
