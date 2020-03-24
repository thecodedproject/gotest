package assert_test

import (
	"github.com/thecodedproject/gotest/assert"
	"testing"
	tfyassert "github.com/stretchr/testify/assert"
)

func TestAssertChannelReceivesWhenPasses(t *testing.T) {

	ch := make(chan interface{}, 1)

	expectedString := "hello world"

	wait := assert.ChannelReceives(t, ch, expectedString)

	ch <- expectedString

	result := wait()
	tfyassert.True(t, result)
}

func TestAssertChannelReceivesWhenFails(t *testing.T) {

	ch := make(chan interface{}, 1)

	expectedString := "hello world"

	fakeT := testing.T{}
	wait := assert.ChannelReceives(&fakeT, ch, expectedString)

	result := wait()
	tfyassert.False(t, result)
}
