package ptr_test

import (
	"testing"

	"github.com/avakarev/go-util/ptr"
	"github.com/avakarev/go-util/testutil"
)

func TestOf(t *testing.T) {
	testutil.Diff(true, *ptr.Of(true), t)
	testutil.Diff(int64(1), *ptr.Of(int64(1)), t)
	testutil.Diff(float64(1.1), *ptr.Of(float64(1.1)), t)
	testutil.Diff("ptr", *ptr.Of("ptr"), t)
}
