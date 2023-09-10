package timeutil_test

import (
	"testing"
	"time"

	"github.com/avakarev/go-util/timeutil"

	"github.com/avakarev/go-util/testutil"
)

func TestClockNow(t *testing.T) {
	std := time.Now().Round(time.Second)
	lib := timeutil.NewClock().Now().Round(time.Second)
	testutil.Diff(std, lib, t)
}

func TestMockNow(t *testing.T) {
	mock := timeutil.NewMock()
	testutil.Diff(time.Unix(0, 0), mock.Now(), t)
	mock.Add(42 * time.Second)
	testutil.Diff(time.Unix(42, 0), mock.Now(), t)
}

func TestMockSince(t *testing.T) {
	mock := timeutil.NewMock()
	start := mock.Now()
	mock.Add(42 * time.Second)
	testutil.Diff(float64(42), mock.Since(start).Seconds(), t)
}

func TestMockUntil(t *testing.T) {
	mock := timeutil.NewMock()
	end := mock.Now().Add(42 * time.Second)
	testutil.Diff(float64(42), mock.Until(end).Seconds(), t)
	mock.Add(2 * time.Second)
	testutil.Diff(float64(40), mock.Until(end).Seconds(), t)
}
