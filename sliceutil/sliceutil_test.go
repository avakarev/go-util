package sliceutil_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/avakarev/go-util/sliceutil"
	"github.com/avakarev/go-util/testutil"
)

func TestShuffle(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e", "f"}
	shuffled := sliceutil.Shuffle(slice)
	testutil.Diff(false, reflect.DeepEqual(slice, shuffled), t)
	slices.Sort(slice)
	slices.Sort(shuffled)
	testutil.Diff(true, reflect.DeepEqual(slice, shuffled), t)
}
