package testutil

import (
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
)

// Diff fails the test if `want` differs from `got`, and prints human-readable error
func Diff(want interface{}, got interface{}, t *testing.T) {
	// ignore all unexported fields
	opts := cmp.FilterPath(func(p cmp.Path) bool {
		sf, ok := p.Index(-1).(cmp.StructField)
		if !ok {
			return false
		}
		r, _ := utf8.DecodeRuneInString(sf.Name())
		return !unicode.IsUpper(r)
	}, cmp.Ignore())

	if diff := cmp.Diff(want, got, opts); diff != "" {
		t.Errorf("%sGot unexpected result (-want +got):\n%s", caller(), diff)
	}
}
