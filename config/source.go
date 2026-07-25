package config

import (
	"github.com/spf13/cast"
)

type Source func() (any, bool)

func DefaultSource(v any) Source {
	return func() (any, bool) {
		return v, true
	}
}

func Lookup[T cast.Basic](ss ...Source) (T, bool) {
	var zero T
	for _, s := range ss {
		v, ok := s()
		if ok {
			if vv, err := cast.ToE[T](v); err == nil {
				return vv, true
			}
		}
	}
	return zero, false
}

func MustLookup[T cast.Basic](ss ...Source) T {
	v, found := Lookup[T](ss...)
	if found {
		return v
	}
	panic("config - not found source or invalid source")
}
