package assert

import (
	"fmt"
  "github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"sort"
)

func LogicallyEqual(
	t *testing.T,
	a any,
	b any,
	s ...any,
) bool {

	if a == nil || b == nil {
		return assert.Equal(t, a, b, s...)
	}

	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)

	if aType != bType {
		return assert.Equal(t, a, b, s...)
	}

	if !valuesLogicallyEqual(
		t,
		reflect.ValueOf(a),
		reflect.ValueOf(b),
		s...,
	) {
		return assert.Equal(t, a, b, s...)
	}

	return true
}

func valuesLogicallyEqual(
	t *testing.T,
	a reflect.Value,
	b reflect.Value,
	s ...any,
) bool {

	if res, ok := maybeCallCmp(a, b); ok {
		return res==0
	}

	if res, ok := maybeCallEqual(a, b); ok {
		return res
	}

	switch a.Kind() {
	case reflect.Ptr:
		return ptrsLogicallyEqual(t, a, b, s...)
	case reflect.Struct:
		return structsLogicallyEqual(t, a, b, s...)
	case reflect.Map:
		return mapsLogicallyEqual(t, a, b, s...)
	case reflect.Slice:
		return slicesLogicallyEqual(t, a, b, s...)
	default:
		return false
	}
}

// maybeCallCmp performs a runtime reflection to see if the type `a` has the
// method `Cmp(rhs TypeOf(b)) int` and calls it if it exists.
func maybeCallCmp(a, b reflect.Value) (cmpResult int64, hasCmp bool) {

	eq := a.MethodByName("Cmp")

	if !eq.IsValid() {
		return 0, false
	}

	if eq.Type().NumIn() != 1 || eq.Type().NumOut() != 1 {
		return 0, false
	}

	if eq.Type().In(0) != b.Type() || eq.Type().Out(0) != reflect.TypeOf(int(0)) {
		return 0, false
	}

	res := eq.Call([]reflect.Value{
		b,
	})

	return res[0].Int(), true
}

// maybeCallEqual performs a runtime reflection to see if the type `a` has the
// method `Equal(rhs TypeOf(b)) bool` and calls it if it exists.
func maybeCallEqual(a, b reflect.Value) (eqResult bool, hasCmp bool) {

	eq := a.MethodByName("Equal")

	if !eq.IsValid() {
		return false, false
	}

	if eq.Type().NumIn() != 1 || eq.Type().NumOut() != 1 {
		return false, false
	}

	if eq.Type().In(0) != b.Type() || eq.Type().Out(0) != reflect.TypeOf(false) {
		return false, false
	}

	res := eq.Call([]reflect.Value{
		b,
	})

	return res[0].Bool(), true
}

func ptrsLogicallyEqual(
	t *testing.T,
	a reflect.Value,
	b reflect.Value,
	s ...any,
) bool {

	if a.IsZero() && b.IsZero() {
		return true
	}

	if a.IsZero() != b.IsZero() {
		return false
	}

	return valuesLogicallyEqual(
		t,
		a.Elem(),
		b.Elem(),
		s...,
	)
}

func structsLogicallyEqual(
	t *testing.T,
	a reflect.Value,
	b reflect.Value,
	s ...any,
) bool {

	retVal := true
	for i:=0; i<a.Type().NumField(); i++ {
		fieldName := a.Type().Field(i).Name
		aField := a.Field(i)
		bField := b.Field(i)

		messageAndFieldName := append(s, "."+fieldName)
		retVal = retVal && valuesLogicallyEqual(t, aField, bField, messageAndFieldName...)
	}

	return retVal
}

func mapsLogicallyEqual(
	t *testing.T,
	a reflect.Value,
	b reflect.Value,
	s ...any,
) bool {

	keysMsg := append([]interface{}{"Keys of map"}, s...)
	ok := assert.Equal(
		t,
		sortedMapKeys(a),
		sortedMapKeys(b),
		keysMsg...,
	)
	if !ok {
		return false
	}

	retval := true
	for _, key := range a.MapKeys() {
		messageAndFieldName := append(s, ".['"+key.String()+"']")

		retval = retval && valuesLogicallyEqual(
			t,
			a.MapIndex(key),
			b.MapIndex(key),
			messageAndFieldName...,
		)
	}

	return retval
}

func slicesLogicallyEqual(
	t *testing.T,
	a reflect.Value,
	b reflect.Value,
	s ...any,
) bool {

	keysMsg := append([]interface{}{"Length of slice"}, s...)
	ok := assert.Equal(
		t,
		a.Len(),
		b.Len(),
		keysMsg...,
	)
	if !ok {
		return false
	}

	retval := true
	for i:=0; i<a.Len(); i++ {
		messageAndFieldName := append(s, fmt.Sprintf(".[%d]", i))

		retval = retval && valuesLogicallyEqual(
			t,
			a.Index(i),
			b.Index(i),
			messageAndFieldName...,
		)
	}

	return retval
}

func sortedMapKeys(value reflect.Value) []string {

	mapKeys := value.MapKeys()
	mapKeysStr := make([]string, 0, len(mapKeys))
	for _, keyVal := range mapKeys {
		mapKeysStr = append(mapKeysStr, keyVal.String())
	}
	sort.Strings(mapKeysStr)
	return mapKeysStr
}
