package assert

import (
	"fmt"
  "github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"sort"
)

func LogicallyEqual(t *testing.T, a, b interface{}, s ...interface{}) bool {

	if a == nil || b == nil {
		return assert.Equal(t, a, b, s...)
	}

	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)

	if aType != bType {
		return assert.Equal(t, a, b, s...)
	}

	if res, ok := maybeCallCmp(a, b); ok {
		message := fmt.Sprint(s...)
		return assert.Equal(
			t,
			int64(0),
			res,
			"%s: Not equal (Cmp method).\n\tExpected: %s\n\tActual: %s",
			message,
			a,
			b,
		)
	}

	if res, ok := maybeCallEqual(a, b); ok {
		message := fmt.Sprint(s...)
		return assert.True(
			t,
			res,
			"%s: Not equal (Equal method).\n\tExpected: %s\n\tActual: %s",
			message,
			a,
			b,
		)
	}

	switch aType.Kind() {
	case reflect.Ptr:
		return ptrsLogicallyEqual(t, a, b, s...)
	case reflect.Struct:
		return structsLogicallyEqual(t, a, b, s...)
	case reflect.Map:
		return mapsLogicallyEqual(t, a, b, s...)
	case reflect.Slice:
		return slicesLogicallyEqual(t, a, b, s...)
	default:
		return assert.Equal(t, a, b, s...)
	}
}

// maybeCallCmp performs a runtime reflection to see if the type `a` has the
// method `Cmp(rhs TypeOf(b)) int` and calls it if it exists.
func maybeCallCmp(a, b any) (cmpResult int64, hasCmp bool) {

	eq := reflect.ValueOf(a).MethodByName("Cmp")

	if !eq.IsValid() {
		return 0, false
	}

	if eq.Type().NumIn() != 1 || eq.Type().NumOut() != 1 {
		return 0, false
	}

	if eq.Type().In(0) != reflect.TypeOf(b) || eq.Type().Out(0) != reflect.TypeOf(int(0)) {
		return 0, false
	}

	res := eq.Call([]reflect.Value{
		reflect.ValueOf(b),
	})

	return res[0].Int(), true
}

// maybeCallEqual performs a runtime reflection to see if the type `a` has the
// method `Equal(rhs TypeOf(b)) bool` and calls it if it exists.
func maybeCallEqual(a, b any) (eqResult bool, hasCmp bool) {

	eq := reflect.ValueOf(a).MethodByName("Equal")

	if !eq.IsValid() {
		return false, false
	}

	if eq.Type().NumIn() != 1 || eq.Type().NumOut() != 1 {
		return false, false
	}

	if eq.Type().In(0) != reflect.TypeOf(b) || eq.Type().Out(0) != reflect.TypeOf(false) {
		return false, false
	}

	res := eq.Call([]reflect.Value{
		reflect.ValueOf(b),
	})

	return res[0].Bool(), true
}

func ptrsLogicallyEqual(
	t *testing.T,
	a interface{},
	b interface{},
	s ...interface{},
) bool {

	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	if aValue.IsZero() || bValue.IsZero() {
		return assert.Equal(t, a, b, s...)
	}

	return LogicallyEqual(
		t,
		aValue.Elem().Interface(),
		bValue.Elem().Interface(),
		s...,
	)
}

func structsLogicallyEqual(
	t *testing.T,
	a interface{},
	b interface{},
	s ...interface{},
) bool {

	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	retVal := true
	publicFields := 0
	for i:=0; i<aValue.Type().NumField(); i++ {
		if aValue.Field(i).CanInterface() {
			fieldName := aValue.Type().Field(i).Name
			aField := aValue.Field(i).Interface()
			bField := bValue.Field(i).Interface()

			messageAndFieldName := append(s, "."+fieldName)
			retVal = retVal && LogicallyEqual(t, aField, bField, messageAndFieldName...)
			publicFields++
		}
	}

	if publicFields > 0 {
		return retVal
	}

	return assert.Equal(t, a, b, s...)
}

func mapsLogicallyEqual(
	t *testing.T,
	a interface{},
	b interface{},
	s ...interface{},
) bool {

	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	keysMsg := append([]interface{}{"Keys of map"}, s...)
	ok := assert.Equal(
		t,
		sortedMapKeys(aValue),
		sortedMapKeys(bValue),
		keysMsg...,
	)
	if !ok {
		return false
	}

	retval := true
	for _, key := range aValue.MapKeys() {
		messageAndFieldName := append(s, ".['"+key.String()+"']")

		retval = retval && LogicallyEqual(
			t,
			aValue.MapIndex(key).Interface(),
			bValue.MapIndex(key).Interface(),
			messageAndFieldName...,
		)
	}

	return retval
}

func slicesLogicallyEqual(
	t *testing.T,
	a interface{},
	b interface{},
	s ...interface{},
) bool {

	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	keysMsg := append([]interface{}{"Length of slice"}, s...)
	ok := assert.Equal(
		t,
		aValue.Len(),
		bValue.Len(),
		keysMsg...,
	)
	if !ok {
		return false
	}

	retval := true
	for i:=0; i<aValue.Len(); i++ {
		messageAndFieldName := append(s, fmt.Sprintf(".[%d]", i))

		retval = retval && LogicallyEqual(
			t,
			aValue.Index(i).Interface(),
			bValue.Index(i).Interface(),
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
