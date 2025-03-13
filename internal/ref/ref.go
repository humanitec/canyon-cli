package ref

import (
	"reflect"
)

func Ref[x any](xi x) *x {
	return &xi
}

func Deref[x any](xi *x, def x) x {
	if xi == nil {
		return def
	}
	return *xi
}

func Coalesce[x any](items ...x) x {
	for _, item := range items {
		if !reflect.ValueOf(item).IsZero() {
			return item
		}
	}
	panic("cannot coalesce with no non-nil items")
}
