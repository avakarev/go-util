package envutil_test

import (
	"errors"
	"testing"
	"time"

	"github.com/avakarev/go-util/envutil"
	"github.com/avakarev/go-util/testutil"
)

func TestStr(t *testing.T) {
	testutil.Diff("", envutil.Str("FOOBAR"), t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	defer resetEnv()
	testutil.Diff("baz qux", envutil.Str("FOOBAR"), t)
}

func TestShouldStr(t *testing.T) {
	empty, err := envutil.ShouldStr("FOOBAR")
	testutil.MustErr(errors.New("env variable \"FOOBAR\" is required, but not set"), err, t)
	testutil.Diff("", empty, t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": ""})
	blank, err := envutil.ShouldStr("FOOBAR")
	testutil.MustErr(errors.New("env variable \"FOOBAR\" is required, but empty"), err, t)
	testutil.Diff("", blank, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	defer resetEnv()
	present, err := envutil.ShouldStr("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff("baz qux", present, t)
}

func TestStrDef(t *testing.T) {
	testutil.Diff("fallback", envutil.StrDef("FOOBAR", "fallback"), t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": ""})
	testutil.Diff("fallback", envutil.StrDef("FOOBAR", "fallback"), t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	defer resetEnv()
	testutil.Diff("baz qux", envutil.StrDef("FOOBAR", "fallback"), t)
}

func TestMustStr(t *testing.T) {
	func() {
		defer func() {
			testutil.Diff(true, recover() != nil, t)
		}()
		_ = envutil.MustStr("FOOBAR")
	}()

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	defer resetEnv()
	testutil.Diff("baz qux", envutil.MustStr("FOOBAR"), t)
}

func TestShouldStrSlice(t *testing.T) {
	_, err := envutil.ShouldStrSlice("FOOBAR", ",")
	testutil.MustErr(errors.New("env variable \"FOOBAR\" is required, but not set"), err, t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz,qux"})
	defer resetEnv()
	value, err := envutil.ShouldStrSlice("FOOBAR", ",")
	testutil.MustNoErr(err, t)
	testutil.Diff([]string{"baz", "qux"}, value, t)
}

func TestMustStrSlice(t *testing.T) {
	func() {
		defer func() {
			testutil.Diff(true, recover() != nil, t)
		}()
		_ = envutil.MustStrSlice("FOOBAR", ",")
	}()

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz,qux"})
	defer resetEnv()
	testutil.Diff([]string{"baz", "qux"}, envutil.MustStrSlice("FOOBAR", ","), t)
}

func TestShouldInt(t *testing.T) {
	num, err := envutil.ShouldInt("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff(int64(0), num, t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	num, err = envutil.ShouldInt("FOOBAR")
	testutil.MustErr(errors.New("strconv.ParseInt: parsing \"baz qux\": invalid syntax"), err, t)
	testutil.Diff(int64(0), num, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "42"})
	num, err = envutil.ShouldInt("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff(int64(42), num, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "-1"})
	num, err = envutil.ShouldInt("FOOBAR")
	testutil.MustNoErr(err, t)
	testutil.Diff(int64(-1), num, t)
	resetEnv()
}

func TestMustInt(t *testing.T) {
	// unset yields 0 without panic
	testutil.Diff(int64(0), envutil.MustInt("FOOBAR"), t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": "42"})
	testutil.Diff(int64(42), envutil.MustInt("FOOBAR"), t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "baz qux"})
	defer resetEnv()
	func() {
		defer func() {
			testutil.Diff(true, recover() != nil, t)
		}()
		_ = envutil.MustInt("FOOBAR")
	}()
}

func TestDurDef(t *testing.T) {
	d, err := envutil.DurDef("FOOBAR", 5*time.Second)
	testutil.MustNoErr(err, t)
	testutil.Diff(5*time.Second, d, t)

	resetEnv := testutil.SetEnv(testutil.Env{"FOOBAR": ""})
	d, err = envutil.DurDef("FOOBAR", 5*time.Second)
	testutil.MustNoErr(err, t)
	testutil.Diff(5*time.Second, d, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "250ms"})
	d, err = envutil.DurDef("FOOBAR", 5*time.Second)
	testutil.MustNoErr(err, t)
	testutil.Diff(250*time.Millisecond, d, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"FOOBAR": "nope"})
	defer resetEnv()
	d, err = envutil.DurDef("FOOBAR", 5*time.Second)
	testutil.MustErr(errors.New("env variable \"FOOBAR\" is invalid duration: time: invalid duration \"nope\""), err, t)
	testutil.Diff(time.Duration(0), d, t)
}
