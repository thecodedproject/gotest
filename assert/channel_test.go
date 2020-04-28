package assert_test

import (
	"github.com/thecodedproject/gotest/assert"
	"testing"
	tfyassert "github.com/stretchr/testify/assert"
)

func TestAssertChannelReceivesOnceWhenPasses(t *testing.T) {

	ch := make(chan interface{}, 1)

	expectedString := "hello world"

	wait := assert.ChannelReceivesOnce(t, ch, expectedString)

	ch <- expectedString

	result := wait()
	tfyassert.True(t, result)
}

func TestAssertChannelReceivesOnceWhenFails(t *testing.T) {

	ch := make(chan interface{}, 1)

	expectedString := "hello world"

	fakeT := testing.T{}
	wait := assert.ChannelReceivesOnce(&fakeT, ch, expectedString)

	result := wait()
	tfyassert.False(t, result)
}

func TestAssertChannelReceivesWhenPasses(t *testing.T) {

	ch := make(chan interface{}, 1)

	expectedStrings := []interface{}{
		"hello",
		"world",
		"!",
	}

	wait := assert.ChannelReceives(t, ch, expectedStrings)

	for _, s := range expectedStrings {
		ch <- s
	}

	result := wait()
	tfyassert.True(t, result)
}

func TestAssertChannelReceivesWhenFails(t *testing.T) {

	ch := make(chan interface{}, 1)

	expectedStrings := []interface{}{
		"hello",
		"world",
		"!",
	}

	fakeT := testing.T{}
	wait := assert.ChannelReceives(&fakeT, ch, expectedStrings)

	ch <- "hello"
	ch <- "!"

	result := wait()
	tfyassert.False(t, result)
}
