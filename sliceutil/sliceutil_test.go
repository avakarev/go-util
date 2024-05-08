package sliceutil_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/avakarev/go-util/sliceutil"
	"github.com/avakarev/go-util/testutil"
)

func TestContainsWithString(t *testing.T) {
	testutil.Diff(false, sliceutil.Contains([]string{"a", "b"}, "c"), t)
	testutil.Diff(true, sliceutil.Contains([]string{"a", "b"}, "b"), t)
}

func TestContainsWithInt(t *testing.T) {
	testutil.Diff(false, sliceutil.Contains([]int{9, 42}, 0), t)
	testutil.Diff(true, sliceutil.Contains([]int{9, 42}, 42), t)
}

func TestContainsWithInt64(t *testing.T) {
	testutil.Diff(false, sliceutil.Contains([]int64{9, 42}, 0), t)
	testutil.Diff(true, sliceutil.Contains([]int64{9, 42}, 42), t)
}

func TestContainsWithFloat64(t *testing.T) {
	testutil.Diff(false, sliceutil.Contains([]float64{9.9, 42.42}, 0.1), t)
	testutil.Diff(true, sliceutil.Contains([]float64{9.9, 42.42}, 42.42), t)
}

func TestShuffle(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e", "f"}
	shuffled := sliceutil.Shuffle(slice)
	testutil.Diff(false, reflect.DeepEqual(slice, shuffled), t)
	slices.Sort(slice)
	slices.Sort(shuffled)
	testutil.Diff(true, reflect.DeepEqual(slice, shuffled), t)
}
