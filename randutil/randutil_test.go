package randutil_test

import (
	"testing"

	"github.com/avakarev/go-util/randutil"
	"github.com/avakarev/go-util/testutil"
)

func TestString(t *testing.T) {
	res, err := randutil.String(42)
	testutil.MustNoErr(err, t)
	testutil.Diff(42, len(res), t)
}

func TestInt(t *testing.T) {
	n := randutil.Int(42)
	testutil.Diff(true, n >= 0 && n <= 42, t)
}
