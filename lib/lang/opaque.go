package lang

import (
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/exp/maps"
)

var (
	cmpOptIgnorePrivate = ignoreUnexportedAlways()
	cmpOptNilIsEmpty    = cmpopts.EquateEmpty()
)

// ignoreUnexportedAlways is a derivative of go-cmp.IgnoreUnexported, but this one
// will always ignore unexported types, recursively.
func ignoreUnexportedAlways() cmp.Option {
	return cmp.FilterPath(
		func(p cmp.Path) bool {
			sf, ok := p.Index(-1).(cmp.StructField)
			if !ok {
				return false
			}
			r, _ := utf8.DecodeRuneInString(sf.Name())
			return !unicode.IsUpper(r)
		},
		cmp.Ignore(),
	)
}

// OpaqueMapsEqual compare maps[<comparable>]<any> for equality, but safely by
// using the cmp package and ignoring un-exported types, and by treating nil/empty
// slices and maps as equal.
func OpaqueMapsEqual[M ~map[K]V, K comparable, V any](m1, m2 M) bool {
	return maps.EqualFunc(m1, m2, func(a, b V) bool {
		return cmp.Equal(a, b,
			cmpOptIgnorePrivate, // ignore all private fields
			cmpOptNilIsEmpty,    // nil/empty slices treated as equal
		)
	})
}
