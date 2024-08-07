package numutil_test

import (
	"testing"

	"github.com/avakarev/go-util/numutil"
	"github.com/avakarev/go-util/testutil"
)

func TestPercent(t *testing.T) {
	testutil.Diff(1.0, numutil.Percent(2, 50), t)
	testutil.Diff(1.25, numutil.Percent(2.5, 50.0), t)
	testutil.Diff(33.33, numutil.Percent(66.66, 50.0), t)
	testutil.Diff(49.0, numutil.Percent(98.0, 50.0), t)
	testutil.Diff(55.55, numutil.Percent(111.1, 50.0), t)
}

func TestPercentOf(t *testing.T) {
	testutil.Diff(2.0, numutil.PercentOf(1, 50), t)
	testutil.Diff(2.5, numutil.PercentOf(1.25, 50.0), t)
	testutil.Diff(66.66, numutil.PercentOf(33.33, 50.0), t)
	testutil.Diff(98.0, numutil.PercentOf(49, 50), t)
	testutil.Diff(111.1, numutil.PercentOf(55.55, 50.0), t)
}

func TestPercentDiff(t *testing.T) {
	testutil.Diff(25.0, numutil.PercentDiff(1.0, 1.25), t)
	testutil.Diff(-20.0, numutil.PercentDiff(1.25, 1.0), t)
}

func TestChangeByPercent(t *testing.T) {
	testutil.Diff(1.25, numutil.ChangeByPercent(1, 25), t)
	testutil.Diff(0.75, numutil.ChangeByPercent(1, -25), t)
	testutil.Diff(26.25, numutil.ChangeByPercent(25, 5), t)
	testutil.Diff(40.0, numutil.ChangeByPercent(50, -20), t)
}

func TestRound(t *testing.T) {
	cases := []struct {
		got       float64
		precision uint
		want      float64
	}{
		{got: 12.3456789, precision: 0, want: 12},
		{got: 12.3456789, precision: 1, want: 12.3},
		{got: 12.3456789, precision: 2, want: 12.35},
		{got: 12.3456789, precision: 3, want: 12.346},
		{got: 12.3456789, precision: 4, want: 12.3457},
		{got: 12.3456789, precision: 5, want: 12.34568},
		{got: 12.3456789, precision: 6, want: 12.345679},
		{got: -12.3456789, precision: 4, want: -12.3457},
	}
	for i := range cases {
		got := numutil.Round(cases[i].got, cases[i].precision)
		testutil.Diff(cases[i].want, got, t)
	}
}

func TestDrawdown(t *testing.T) {
	testutil.Diff(56.7435, numutil.Round(numutil.Drawdown(40.78, 17.64), 4), t)
	testutil.Diff(54.9534, numutil.Round(numutil.Drawdown(40.78, 18.37), 4), t)
	testutil.Diff(50.0, numutil.Drawdown(90, 45), t)
	testutil.Diff(1.0, numutil.Drawdown(100, 99), t)
	testutil.Diff(-25.0, numutil.Drawdown(20, 25), t)
}
