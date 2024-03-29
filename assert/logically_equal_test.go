package assert_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
	tfyassert "github.com/stretchr/testify/assert"

	"github.com/thecodedproject/gotest/assert"
)

type MyCmp struct {
	A int
	Unused int
}

func (m MyCmp) Cmp(rhs MyCmp) int {

	if m.A > rhs.A {
		return 1
	}
	if m.A < rhs.A {
		return -1
	}
	return 0
}

type MyEqual struct {
	A int
	Unused int
}

func (m MyEqual) Equal(rhs MyEqual) bool {
	return m.A == rhs.A
}

type MyStruct struct {
	Exported string
	unexported string
}

type MyNestedStruct struct {
	ExpNest MyStruct
}

type NestedEqual struct {
	
}

func TestBasicTypes(t *testing.T) {
	testLogicallyEqual(t, true, true, true)
	testLogicallyEqual(t, int64(10), int64(10), true)
	testLogicallyEqual(t, 1.234, 1.234, true)
	testLogicallyEqual(t, "abc", "abc", true)
	testLogicallyEqual(t, int64(10), int64(11), false)
	testLogicallyEqual(t, int64(10), "hello", false)
	testLogicallyEqual(t, "abc", false, false)
}

func TestStructs(t *testing.T) {

	testLogicallyEqualWithDesc(
		t,
		"all fields equal",
		MyStruct{Exported: "a", unexported: "b"},
		MyStruct{Exported: "a", unexported: "b"},
		true,
	)
	testLogicallyEqualWithDesc(
		t,
		"unexported fields not equal",
		MyStruct{Exported: "a", unexported: "b"},
		MyStruct{Exported: "a", unexported: "c"},
		false,
	)

	testLogicallyEqualWithDesc(
		t,
		"Cmp method compares equal",
		MyCmp{A: 10, Unused: 1},
		MyCmp{A: 10, Unused: 2},
		true,
	)

	testLogicallyEqualWithDesc(
		t,
		"Cmp method compares not equal",
		MyCmp{A: 10, Unused: 1},
		MyCmp{A: 11, Unused: 1},
		false,
	)

	testLogicallyEqualWithDesc(
		t,
		"Equal method compares equal",
		MyEqual{A: 20, Unused: 1},
		MyEqual{A: 20, Unused: 2},
		true,
	)

	testLogicallyEqualWithDesc(
		t,
		"Equal method compares not equal",
		MyEqual{A: 20, Unused: 1},
		MyEqual{A: 21, Unused: 1},
		false,
	)

	testLogicallyEqualWithDesc(
		t,
		"Uninitialised nested elements which are equal",
		MyNestedStruct{},
		MyNestedStruct{},
		true,
	)

	testLogicallyEqualWithDesc(
		t,
		"One uninitialised nested struct compares not equal",
		MyNestedStruct{ExpNest: MyStruct{Exported: "abc"}},
		MyNestedStruct{},
		false,
	)

	testLogicallyEqualWithDesc(
		t,
		"Nested structs equal",
		MyNestedStruct{ExpNest: MyStruct{Exported: "abc"}},
		MyNestedStruct{ExpNest: MyStruct{Exported: "abc"}},
		true,
	)
}

func testLogicallyEqual[A any, B any](
	t *testing.T,
	a A,
	b B,
	expectedResult bool,
) {
	testLogicallyEqualWithDesc(
		t,
		fmt.Sprintf("%t", expectedResult),
		a,
		b,
		expectedResult,
	)
}

func testLogicallyEqualWithDesc[A any, B any](
	t *testing.T,
	description string,
	a A,
	b B,
	expectedResult bool,
) {
	testName := fmt.Sprintf("%s_%s_%s", reflect.TypeOf(a).String(), reflect.TypeOf(b).String(), description)
	t.Run(testName, func(t *testing.T) {
		var fakeT testing.T
		res := assert.LogicallyEqual(&fakeT, a, b)
		tfyassert.Equal(t, expectedResult, res)
	})

	t.Run("pointers_" + testName, func(t *testing.T) {
		var fakeT testing.T
		res := assert.LogicallyEqual(&fakeT, &a, &b)
		tfyassert.Equal(t, expectedResult, res)
	})
}

func TestLogicallyEqual(t *testing.T) {

	testCases := []struct{
		name string
		a interface{}
		b interface{}
		s []interface{}
		pass bool
	}{
		{
			name: "integers",
			a: int(1),
			b: int(1),
			pass: true,
		},
		{
			name: "shopspring decimals equal",
			a: decimal.NewFromFloat(2.0),
			b: decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
			pass: true,
		},
		{
			name: "shopspring decimals not equal",
			a: decimal.NewFromFloat(2.0),
			b: decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10)),
			pass: false,
		},
		{
			name: "custom comparable type equal",
			a: MyCmp{3, 1},
			b: MyCmp{3, 2},
			pass: true,
		},
		{
			name: "custom comparable type not equal",
			a: MyCmp{3, 1},
			b: MyCmp{4, 2},
			pass: false,
		},
		{
			name: "custom equal type equal",
			a: MyEqual{2, 5},
			b: MyEqual{2, 7},
			pass: true,
		},
		{
			name: "custom equal type not equal",
			a: MyEqual{1, 5},
			b: MyEqual{6, 7},
			pass: false,
		},
		{
			name: "nil errors inside structs equal",
			a: struct{
				Err error
			}{},
			b: struct{
				Err error
			}{},
			pass: true,
		},
		{
			name: "nil errors inside structs not equal - one not initalised",
			a: struct{
				Err error
			}{},
			b: struct{
				Err error
			}{
				Err: errors.New("some error"),
			},
			pass: false,
		},
		{
			name: "shopspring decimals inside struct equal",
			a: struct{
				Field decimal.Decimal
			}{decimal.NewFromFloat(2)},
			b: struct{
				Field decimal.Decimal
			}{decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10))},
			pass: true,
		},
		{
			name: "shopspring decimals inside struct not equal",
			a: struct{
				Field decimal.Decimal
			}{decimal.NewFromFloat(2)},
			b: struct{
				Field decimal.Decimal
			}{decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10))},
			pass: false,
		},
		{
			name: "struct with no public fields is compared physically equal when equal",
			a: struct{
				privateOne int
				privateTwo string
			}{1, "a"},
			b: struct{
				privateOne int
				privateTwo string
			}{1, "a"},
			pass: true,
		},
		{
			name: "struct with no public fields is compared physically equal when not equal",
			a: struct{
				privateOne int
				privateTwo string
			}{1, "a"},
			b: struct{
				privateOne int
				privateTwo string
			}{1, "b"},
			pass: false,
		},
		{
			name: "map of decimals when equal",
			a: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(2),
				"two": decimal.NewFromFloat(0),
			},
			b: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
				"two": decimal.Decimal{},
			},
			pass: true,
		},
		{
			name: "map of decimals when not equal",
			a: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(2),
				"two": decimal.NewFromFloat(0),
			},
			b: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10)),
				"two": decimal.Decimal{},
			},
			pass: false,
		},
		{
			name: "map of decimals with different field names",
			a: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(2),
				"two": decimal.NewFromFloat(0),
			},
			b: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
				"three": decimal.Decimal{},
			},
			pass: false,
		},
		{
			name: "map of decimals with different lengths - a contains more entries",
			a: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(2),
				"two": decimal.NewFromFloat(0),
			},
			b: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
			},
			pass: false,
		},
		{
			name: "map of decimals with different lengths - b contains more entries",
			a: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(2),
			},
			b: map[string]decimal.Decimal{
				"one": decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
				"two": decimal.NewFromFloat(0),
			},
			pass: false,
		},
		{
			name: "map inside struct when equal",
			a: struct{
				M map[string]decimal.Decimal
			}{
				M: map[string]decimal.Decimal{
					"one": decimal.NewFromFloat(2),
					"two": decimal.Decimal{},
				},
			},
			b: struct{
				M map[string]decimal.Decimal
			}{
				M: map[string]decimal.Decimal{
					"one": decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
					"two": decimal.NewFromFloat(0),
				},
			},
			pass: true,
		},
		{
			name: "map inside struct when not equal",
			a: struct{
				M map[string]decimal.Decimal
			}{
				M: map[string]decimal.Decimal{
					"one": decimal.NewFromFloat(2),
					"two": decimal.Decimal{},
				},
			},
			b: struct{
				M map[string]decimal.Decimal
			}{
				M: map[string]decimal.Decimal{
					"one": decimal.NewFromFloat(10).Div(decimal.NewFromFloat(10)),
					"two": decimal.NewFromFloat(0),
				},
			},
			pass: false,
		},
		{
			name: "slice of decimals when equal",
			a: []decimal.Decimal{
				decimal.NewFromFloat(2),
				decimal.NewFromFloat(0),
			},
			b: []decimal.Decimal{
				decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
				decimal.Decimal{},
			},
			pass: true,
		},
		{
			name: "slice of decimals when not equal",
			a: []decimal.Decimal{
				decimal.NewFromFloat(2),
				decimal.NewFromFloat(0),
			},
			b: []decimal.Decimal{
				decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10)),
				decimal.Decimal{},
			},
			pass: false,
		},
		{
			name: "slice with different lengths - a contains more than b",
			a: []decimal.Decimal{
				decimal.NewFromFloat(2),
				decimal.NewFromFloat(0),
			},
			b: []decimal.Decimal{
				decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
			},
			pass: false,
		},
		{
			name: "slice with different lengths - b contains more than a",
			a: []decimal.Decimal{
				decimal.NewFromFloat(0),
			},
			b: []decimal.Decimal{
				decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10)),
				decimal.Decimal{},
			},
			pass: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			var fakeT testing.T
			res := assert.LogicallyEqual(&fakeT, test.a, test.b, test.s...)
			tfyassert.Equal(t, test.pass, res)
		})
	}
}

func TestLogicallyEqualWithPtrs(t *testing.T) {

	testCases := []struct{
		name string
		aCtr func() interface{}
		bCtr func() interface{}
		s []interface{}
		pass bool
	}{
		{
			name: "integers equal",
			aCtr: func() interface{} {
				i := 1
				return &i
			},
			bCtr: func() interface{} {
				i := 1
				return &i
			},
			pass: true,
		},
		{
			name: "integers not equal",
			aCtr: func() interface{} {
				i := 1
				return &i
			},
			bCtr: func() interface{} {
				i := 2
				return &i
			},
			pass: false,
		},
		{
			name: "shopspring decimals equal",
			aCtr: func() interface{} {
				i := decimal.NewFromFloat(2.0)
				return &i
			},
			bCtr: func() interface{} {
				i := decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10))
				return &i
			},
			pass: true,
		},
		{
			name: "shopspring decimals not equal",
			aCtr: func() interface{} {
				i := decimal.NewFromFloat(2.0)
				return &i
			},
			bCtr: func() interface{} {
				i := decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10))
				return &i
			},
			pass: false,
		},
		{
			name: "shopspring decimals inside struct equal",
			aCtr: func() interface{} {
				i := decimal.NewFromFloat(2.0)
				return struct{
					Field *decimal.Decimal
				}{
					Field: &i,
				}
			},
			bCtr: func() interface{} {
				i := decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10))
				return struct{
					Field *decimal.Decimal
				}{
					Field: &i,
				}
			},
			pass: true,
		},
		{
			name: "shopspring decimals inside struct not equal",
			aCtr: func() interface{} {
				i := decimal.NewFromFloat(2.0)
				return struct{
					Field *decimal.Decimal
				}{
					Field: &i,
				}
			},
			bCtr: func() interface{} {
				i := decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10))
				return struct{
					Field *decimal.Decimal
				}{
					Field: &i,
				}
			},
			pass: false,
		},
		{
			name: "shopspring decimals inside struct not equal with one being zero value",
			aCtr: func() interface{} {
				i := decimal.NewFromFloat(2.0)
				return struct{
					Field *decimal.Decimal
				}{
					Field: &i,
				}
			},
			bCtr: func() interface{} {
				return struct{
					Field *decimal.Decimal
				}{}
			},
			pass: false,
		},
		{
			name: "shopspring decimals inside struct when both structs are zero value",
			aCtr: func() interface{} {
				return struct{
					Field *decimal.Decimal
				}{}
			},
			bCtr: func() interface{} {
				return struct{
					Field *decimal.Decimal
				}{}
			},
			pass: true,
		},
		{
			name: "shopspring decimals inside struct not equal with one being nil",
			aCtr: func() interface{} {
				i := decimal.NewFromFloat(2.0)
				return struct{
					Field *decimal.Decimal
				}{
					Field: &i,
				}
			},
			bCtr: func() interface{} {
				return struct{
					Field *decimal.Decimal
				}{
					Field: nil,
				}
			},
			pass: false,
		},
		{
			name: "shopspring decimals inside struct equal with both being nil",
			aCtr: func() interface{} {
				return struct{
					Field *decimal.Decimal
				}{
					Field: nil,
				}
			},
			bCtr: func() interface{} {
				return struct{
					Field *decimal.Decimal
				}{
					Field: nil,
				}
			},
			pass: true,
		},
		{
			name: "map of decimals when equal",
			aCtr: func() interface{} {
				one := decimal.NewFromFloat(2.0)
				two := decimal.NewFromFloat(0)
				return map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
			},
			bCtr: func() interface{} {
				one := decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10))
				two := decimal.Decimal{}
				return map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
			},
			pass: true,
		},
		{
			name: "map of decimals when not equal",
			aCtr: func() interface{} {
				one := decimal.NewFromFloat(2.0)
				two := decimal.NewFromFloat(0)
				return map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
			},
			bCtr: func() interface{} {
				one := decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10))
				two := decimal.Decimal{}
				return map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
			},
			pass: false,
		},
		{
			name: "map inside struct when equal",
			aCtr: func() interface{} {
				one := decimal.NewFromFloat(2.0)
				two := decimal.NewFromFloat(0)
				m := map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
				return struct {
					M *map[string]*decimal.Decimal
				}{
					M: &m,
				}
			},
			bCtr: func() interface{} {
				one := decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10))
				two := decimal.Decimal{}
				m := map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
				return struct {
					M *map[string]*decimal.Decimal
				}{
					M: &m,
				}
			},
			pass: true,
		},
		{
			name: "map inside struct when not equal",
			aCtr: func() interface{} {
				one := decimal.NewFromFloat(2.0)
				two := decimal.NewFromFloat(0)
				m := map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
				return struct {
					M *map[string]*decimal.Decimal
				}{
					M: &m,
				}
			},
			bCtr: func() interface{} {
				one := decimal.NewFromFloat(40).Div(decimal.NewFromFloat(10))
				two := decimal.Decimal{}
				m := map[string]*decimal.Decimal{
					"one": &one,
					"two": &two,
				}
				return struct {
					M *map[string]*decimal.Decimal
				}{
					M: &m,
				}
			},
			pass: false,
		},
		{
			name: "slice of decimals when equal",
			aCtr: func() interface{} {
				one := decimal.NewFromFloat(2.0)
				two := decimal.NewFromFloat(0)
				return []*decimal.Decimal{&one, &two}
			},
			bCtr: func() interface{} {
				one := decimal.NewFromFloat(20).Div(decimal.NewFromFloat(10))
				two := decimal.Decimal{}
				return []*decimal.Decimal{&one, &two}
			},
			pass: true,
		},
		{
			name: "slice of decimals when not equal",
			aCtr: func() interface{} {
				one := decimal.NewFromFloat(2.0)
				two := decimal.NewFromFloat(0)
				return []*decimal.Decimal{&one, &two}
			},
			bCtr: func() interface{} {
				one := decimal.NewFromFloat(30).Div(decimal.NewFromFloat(10))
				two := decimal.Decimal{}
				return []*decimal.Decimal{&one, &two}
			},
			pass: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			var fakeT testing.T
			a := test.aCtr()
			b := test.bCtr()
			res := assert.LogicallyEqual(&fakeT, a, b, test.s...)
			tfyassert.Equal(t, test.pass, res)
		})
	}
}
