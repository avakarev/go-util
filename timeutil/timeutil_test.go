package timeutil_test

import (
	"testing"
	"time"

	"github.com/avakarev/go-util/timeutil"

	"github.com/avakarev/go-util/testutil"
)

func timeFixture() time.Time {
	t, _ := time.Parse(time.RFC3339, "2022-06-03T16:26:15Z")
	return t
}

func TestLocalWithNoTZ(t *testing.T) {
	reset := testutil.SetEnv(testutil.Env{"TZ": ""})
	defer reset()

	testutil.MustNoErr(timeutil.Init(), t)
	testutil.Diff("2022-06-03T16:26:15Z", timeutil.Local(timeFixture()).Format(time.RFC3339), t)
}

func TestLocalWithCET(t *testing.T) {
	reset := testutil.SetEnv(testutil.Env{"TZ": "Europe/Berlin"})
	defer reset()

	testutil.MustNoErr(timeutil.Init(), t)
	testutil.Diff("2022-06-03T18:26:15+02:00", timeutil.Local(timeFixture()).Format(time.RFC3339), t)
}

func TestLocalWithEET(t *testing.T) {
	reset := testutil.SetEnv(testutil.Env{"TZ": "Europe/Bucharest"})
	defer reset()

	testutil.MustNoErr(timeutil.Init(), t)
	testutil.Diff("2022-06-03T19:26:15+03:00", timeutil.Local(timeFixture()).Format(time.RFC3339), t)
}

func TestMockNowFn(t *testing.T) {
	timeutil.MockNow(timeFixture)
	defer timeutil.UnmockNow()

	testutil.Diff(timeFixture(), timeutil.Now(), t)
}
