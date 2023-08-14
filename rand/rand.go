package rand

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"

)

var (
	maxContainerSize = 5
)

func New[Type any](t testing.TB) Type {
	return NewFromSeed[Type](t, time.Now().UnixNano())
}

func NewFromSeed[Type any](t testing.TB, seed int64) Type {
	var toFill Type
	FillFromSeed(t, &toFill, seed)
	return toFill
}

func Fill(t testing.TB, toFill any) {
	FillFromSeed(t, toFill, time.Now().UnixNano())
}

func FillFromSeed(t testing.TB, toFill any, seed int64) {
	r := rand.New(rand.NewSource(seed))
	v := reflect.ValueOf(toFill)
	fillValue(t, v, r)
}

func SetMaxContainerSize(t testing.TB, n int) {
	if n < 0 {
		require.Fail(t, "max container size cannot be less than 0", n)
	}
	maxContainerSize = n
}

func fillValue(t testing.TB, v reflect.Value, r *rand.Rand) {

	if v.Kind() != reflect.Pointer && !v.CanAddr() {
		require.Fail(t, "gotest/rand: cannot fill unaddressable value - value should be passed by reference")
	}

	switch v.Kind() {
	case reflect.Array:
		n := v.Len()
		for i:=0; i<n; i++ {
			fillValue(t, v.Index(i), r)
		}
	case reflect.Bool:
		if v.CanSet() {
			v.SetBool(r.Int()%2 == 0)
		}
	case reflect.Complex64:
		if v.CanSet() {
			v.SetComplex(complex(
				r.Float64(),
				r.Float64(),
			))
		}
	case reflect.Complex128:
		if v.CanSet() {
			v.SetComplex(complex(
				r.Float64(),
				r.Float64(),
			))
		}
	case reflect.Float32:
		if v.CanSet() {
			v.SetFloat(r.Float64())
		}
	case reflect.Float64:
		if v.CanSet() {
			v.SetFloat(r.Float64())
		}
	case reflect.Int:
		if v.CanSet() {
			v.SetInt(r.Int63())
		}
	case reflect.Int8:
		if v.CanSet() {
			v.SetInt(r.Int63())
		}
	case reflect.Int16:
		if v.CanSet() {
			v.SetInt(r.Int63())
		}
	case reflect.Int32:
		if v.CanSet() {
			v.SetInt(r.Int63())
		}
	case reflect.Int64:
		if v.CanSet() {
			v.SetInt(r.Int63())
		}
	case reflect.Interface:
		if v.CanInterface() {
			fillValue(t, v.Elem(), r)
		}
	case reflect.Map:
		n := v.Len()
		if v.Len() == 0 {
			n = r.Intn(maxContainerSize - 1) + 1
		}

		v.Set(reflect.MakeMapWithSize(v.Type(), n))

		for i:=0; i<n; i++ {
			k := reflect.New(v.Type().Key())
			fillValue(t, k.Elem(), r)

			val := reflect.New(v.Type().Elem())
			fillValue(t, val.Elem(), r)

			v.SetMapIndex(k.Elem(), val.Elem())
		}
	case reflect.Pointer:
		if v.IsZero() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		fillValue(t, reflect.Indirect(v), r)
	case reflect.Slice:
		if v.Len() != 0 {
			for i:=0; i<v.Len(); i++ {
				fillValue(t, v.Index(i), r)
			}
			return
		}

		if v.Cap() == 0 {
			n := r.Intn(maxContainerSize - 1) + 1
			v.Grow(n)
			// Grow _may_ set the capasicty to something larger than n;
			// therefore we explictly set the capacity as well
			v.SetCap(n)
		}

		for i:=0; i<v.Cap(); i++ {
			e := reflect.New(v.Type().Elem())
			fillValue(t, e, r)
			v.Set(reflect.Append(v, reflect.Indirect(e)))
		}
	case reflect.String:
		if v.CanSet() {
			v.SetString(fmt.Sprintf("%x", r.Uint64()))
		}
	case reflect.Struct:
		n := v.NumField()
		for i:=0; i<n; i++ {
			f := v.Field(i)
			if f.CanSet() {
				fillValue(t, f, r)
			} else {
				newF := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr()))
				fillValue(t, newF.Elem(), r)
			}
		}
	case reflect.Uint:
		if v.CanSet() {
			v.SetUint(r.Uint64())
		}
	case reflect.Uint8:
		if v.CanSet() {
			v.SetUint(r.Uint64())
		}
	case reflect.Uint16:
		if v.CanSet() {
			v.SetUint(r.Uint64())
		}
	case reflect.Uint32:
		if v.CanSet() {
			v.SetUint(r.Uint64())
		}
	case reflect.Uint64:
		if v.CanSet() {
			v.SetUint(r.Uint64())
		}
	case reflect.Uintptr:
		if v.CanSet() {
			v.SetUint(r.Uint64())
		}
	}
}
