package tnglib

import "github.com/d5/tengo/v2"

// Ternary condition in one line, for simple assignment.
func Ternary[T any](cond bool, ifTrue T, ifFalse T) T {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// ToObject simply return argument and true
func ToObject(o tengo.Object) (tengo.Object, bool) {
	return o, true
}

// MapGet return value with given key in case object is map.
// If object is not map or value does not exists, it will return defaultV
func MapGet[T any](o tengo.Object, key string, defVal T, fn func(tengo.Object) (T, bool)) T {
	if o == nil {
		return defVal
	}

	var m map[string]tengo.Object
	switch vm := o.(type) {
	case *tengo.Map:
		m = vm.Value
	case *tengo.ImmutableMap:
		m = vm.Value
	}
	if len(m) == 0 {
		return defVal
	}

	v, ok := m[key]
	if !ok {
		return defVal
	}
	vo, ok := fn(v)
	if !ok {
		return defVal
	}
	return vo
}

// ValueGet return underlying object value or defValue if object is not convertible to T
func ValueGet[T any](o tengo.Object, defVal T, fn func(tengo.Object) (T, bool)) T {
	if o == nil {
		return defVal
	}
	vo, ok := fn(o)
	if !ok {
		return defVal
	}
	return vo
}

// LookupLoop implement switch-case like lookup table functionality.
func LookupLoop[K comparable, V any](key K, keys []K, values []V) (V, bool) {
	nk := len(keys)
	nv := len(values)

	var zeroVal V
	if nk == 0 || nv == 0 || nk != nv {
		return zeroVal, false
	}
	for i := 0; i < nk; i++ {
		if key == keys[i] {
			return values[i], true
		}
	}
	return zeroVal, false
}
