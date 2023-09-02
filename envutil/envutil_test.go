package envutil_test

import (
	"errors"
	"testing"

	"github.com/avakarev/go-util/envutil"
	"github.com/avakarev/go-util/testutil"
)

func TestStr(t *testing.T) {
	testutil.Diff("", envutil.Str("FOOBAR"), t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	defer resetEnv()
	testutil.Diff("baz qux", envutil.Str("FOOBAR"), t)
}

func TestMustStr(t *testing.T) {
	empty, err := envutil.MustStr("FOOBAR")
	testutil.MustErr(errors.New("environment variable \"FOOBAR\" is not set"), err, t)
	testutil.Diff("", empty, t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	defer resetEnv()
	present, err := envutil.MustStr("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff("baz qux", present, t)
}

func TestStrSlice(t *testing.T) {
	testutil.Diff("", envutil.Str("FOOBAR"), t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz,qux"})
	defer resetEnv()
	testutil.Diff([]string{"baz", "qux"}, envutil.StrSlice("FOOBAR", ","), t)
}

func TestInt(t *testing.T) {
	num, err := envutil.Int("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff(int64(0), num, t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	num, err = envutil.Int("FOOBAR")
	testutil.MustErr(errors.New("strconv.ParseInt: parsing \"baz qux\": invalid syntax"), err, t)
	testutil.Diff(int64(0), num, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "42"})
	num, err = envutil.Int("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff(int64(42), num, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "-1"})
	num, err = envutil.Int("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff(int64(-1), num, t)
	resetEnv()
}
