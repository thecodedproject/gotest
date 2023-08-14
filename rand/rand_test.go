package rand_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/thecodedproject/gotest/assert"
	"github.com/thecodedproject/gotest/rand"
)

type MyStruct struct {
	Exported string
	unexported string
}

type MyNestedStruct struct {
	Exported string
	unexported string
	ExpNest MyStruct
	unexpNest MyStruct
}

func TestNew(t *testing.T) {

	t.Run("int64", func(t *testing.T) {
		actual := rand.New[int64](t)
		require.Equal(t, "int64", reflect.TypeOf(actual).String())
	})

	t.Run("string", func(t *testing.T) {
		actual := rand.New[string](t)
		require.Equal(t, "string", reflect.TypeOf(actual).String())
	})

	t.Run("struct", func(t *testing.T) {
		actual := rand.New[MyNestedStruct](t)
		require.Equal(t, "rand_test.MyNestedStruct", reflect.TypeOf(actual).String())
	})
}

func TestFill(t *testing.T) {

	t.Run("int64", func(t *testing.T) {
		var actual int64
		rand.Fill(t, &actual)
		require.Equal(t, "int64", reflect.TypeOf(actual).String())
	})

	t.Run("string", func(t *testing.T) {
		var actual string
		rand.Fill(t, &actual)
		require.Equal(t, "string", reflect.TypeOf(actual).String())
	})

	t.Run("struct", func(t *testing.T) {
		var actual MyNestedStruct
		rand.Fill(t, &actual)
		require.Equal(t, "rand_test.MyNestedStruct", reflect.TypeOf(actual).String())
	})
}

func TestBasicTypes(t *testing.T) {
	testFromSeedWithExpected(t, 1233, bool(false))
	testFromSeedWithExpected(t, 1234, bool(true))
	var c64 complex64
	c64 = 0.7551405+0.33032042i
	testFromSeedWithExpected(t, 2345, c64)
	var c128 complex128
	c128 = 0.7551405130308153+0.3303204161281608i
	testFromSeedWithExpected(t, 2345, c128)
	testFromSeedWithExpected(t, 1134, float32(0.48107496))
	testFromSeedWithExpected(t, 1135, float64(0.00466190836082573))
	testFromSeedWithExpected(t, 1234, int(2041104533947223744))
	testFromSeedWithExpected(t, 1234, int8(-64))
	testFromSeedWithExpected(t, 1234, int16(3776))
	testFromSeedWithExpected(t, 1234, int32(1734151872))
	testFromSeedWithExpected(t, 1234, int64(2041104533947223744))
	testFromSeedWithExpected(t, 1222, string("c2cbc28bb1abb4fe"))
	testFromSeedWithExpected(t, 1234, uint(0x9c5375c2675d0ec0))
	testFromSeedWithExpected(t, 1234, uint16(0xec0))
	testFromSeedWithExpected(t, 1234, uint32(0x675d0ec0))
	testFromSeedWithExpected(t, 1234, uint64(0x9c5375c2675d0ec0))
	testFromSeedWithExpected(t, 1234, uint8(0xc0))
	testFromSeedWithExpected(t, 1235, uintptr(0xe30dd98821827b96))
}

func TestArrays(t *testing.T) {
	testFromSeedWithExpected(t, 1234, [4]int{
		2041104533947223744,
		6276915669504994697,
		6006156956070140861,
		8290995405146611925,
	})
	testFromSeedWithExpected(t, 1234, [2]MyStruct{
		{
			Exported: "9c5375c2675d0ec0",
			unexported: "571c158b7dedad89",
		},
		{
			Exported: "d35a27edf7cd5bbd",
			unexported: "730f8864b628e0d5",
		},
	})
	testFromSeedWithExpected(t, 1234, [2][2]uint8{
		{0xc0,0x89},
		{0xbd,0xd5},
	})
}

func TestSlices(t *testing.T) {

	rand.SetMaxContainerSize(t, 5)

	testFromSeedWithExpected(t, 2345, []int{
		3046668089318711591,
		7744640111461573652,
		966966351434040366,
		4697278898334172992,
	})

	testFromSeedWithExpected(t, 31, []MyStruct{
		{
			Exported: "59db5047520fd39a",
			unexported: "86b3e6428f5f4be",
		},
	})

	t.Run("slice with non-zero capactiy and zero length", func(t *testing.T) {
		toFill := make([]string, 0, 2)
		expected := []string{
			"9eed321a7dd9247b",
			"6447eb3bc4ff7711",
		}
		rand.FillFromSeed(t, &toFill, 454)
		require.Equal(t, expected, toFill)
	})

	t.Run("slice with non-zero length", func(t *testing.T) {
		toFill := make([]bool, 7, 7)
		expected := []bool{true, false, true, true, false, false, true}
		rand.FillFromSeed(t, &toFill, 567)
		require.Equal(t, expected, toFill)
	})

	t.Run("nil slice", func(t *testing.T) {
		var toFill []bool
		require.Nil(t, toFill)
		expected := []bool{false, false, true, false}
		rand.FillFromSeed(t, &toFill, 5678)
		require.Equal(t, expected, toFill)
	})
}

func TestMaps(t *testing.T) {

	rand.SetMaxContainerSize(t, 7)

	testFromSeedWithExpected(t, 2345, map[int8]bool{
		39: true,
		46: true,
	})

	s1 := "6b7a7d2604e8a414"
	s2 := "4130164552258740"
	testFromSeedWithExpected(t, 2345, map[string]*string{
		"aa47f07c3c5c2127": &s1,
		"d6b5ac5fefb442e": &s2,
	})

	testFromSeedWithExpected(t, 28921, map[string]MyStruct{
		"5f437a83712fbeb5": MyStruct{
			Exported: "f78c2f02c8103c6e",
			unexported: "cb605d31c09f9f23",
		},
		"b47e04796f3b50b1": MyStruct{
			Exported: "fb11933b032e5ff0",
			unexported: "7f00fc55423901ae",
		},
	})

	t.Run("map with existing keys", func(t *testing.T) {
		toFill := map[string]string{
			"a": "",
			"b": "",
			"c": "",
			"d": "",
		}
		expected := map[string]string{
			"4aab89e249f74231": "1eff202e48b40288",
			"501734048870ff15": "95921262884a19ae",
			"65ee0562e50f71a7": "ba30969b2e3c7493",
			"b2d8f6bbc3d1a071": "4df4823d2d3124a5",
		}
		rand.FillFromSeed(t, &toFill, 7238)
		require.Equal(t, expected, toFill)
	})

	t.Run("nil map", func(t *testing.T) {
		var toFill map[bool]string
		require.Nil(t, toFill)
		expected := map[bool]string{
			false: "32ab952735d02289",
			true: "7cb30c6c5f6fbbd0",
		}
		rand.FillFromSeed(t, &toFill, 5469)
		require.Equal(t, expected, toFill)
	})

}

func TestStructs(t *testing.T) {
	testFromSeedWithExpected(t, 1235, MyStruct{
		Exported: "e30dd98821827b96",
		unexported: "8b931aa326f28080",
	})

	testFromSeedWithExpected(t, 1122, MyNestedStruct{
		Exported: "e0159c5ee944c625",
		unexported: "3b0d2af7f3adf14e",
		ExpNest: MyStruct{
			Exported: "7789f5afa94e4b8a",
			unexported: "236018edafb9c652",
		},
		unexpNest: MyStruct{
			Exported: "15a5b1c7c90ca169",
			unexported: "ea16685e570c9d18",
		},
	})

	// Time needs to use LogicallyEqual as the monotonic clock cannot be initalised
	testFillFromSeedLogicallyEqual(t, 1235,
		time.Date(2095, time.August, 25, 14, 16, 0, 562199446, time.UTC),
	)
}

func testFromSeedWithExpected[F any](t *testing.T, seed int64, expected F) {
	testName := fmt.Sprintf("%s_%d", reflect.TypeOf(expected).String(), seed)
	t.Run("FillFromSeed_" + testName, func(t *testing.T) {
		var toFill F
		rand.FillFromSeed(t, &toFill, seed)
		require.Equal(t, expected, toFill)
	})

	t.Run("NewFromSeed_" + testName, func(t *testing.T) {
		actual := rand.NewFromSeed[F](t, seed)
		require.Equal(t, expected, actual)
	})
}


func testFillFromSeedLogicallyEqual[F any](t *testing.T, seed int64, expected F) {
	testName := fmt.Sprintf("%s_%d", reflect.TypeOf(expected).String(), seed)
	t.Run(testName, func(t *testing.T) {
		var toFill F
		rand.FillFromSeed(t, &toFill, seed)
		assert.LogicallyEqual(t, expected, toFill)
	})
}
