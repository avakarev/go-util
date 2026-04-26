package timeutil_test

import (
	"testing"
	"time"

	"github.com/avakarev/go-util/timeutil"

	"github.com/avakarev/go-util/testutil"
)

func timeFixture() time.Time {
	return time.Date(2022, 6, 3, 16, 26, 15, 0, time.UTC)
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

func TestStartOfDay(t *testing.T) {
	testutil.Diff("2022-06-03T00:00:01Z", timeutil.StartOfDay(timeFixture()).Format(time.RFC3339), t)
}

func TestEndOfDay(t *testing.T) {
	testutil.Diff("2022-06-03T23:59:59Z", timeutil.EndOfDay(timeFixture()).Format(time.RFC3339), t)
}
