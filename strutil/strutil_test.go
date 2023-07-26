package strutil_test

import (
	"testing"

	"github.com/avakarev/go-util/strutil"
	"github.com/avakarev/go-util/testutil"
)

func TestMaskRight(t *testing.T) {
	cases := []struct {
		s    string
		n    int
		want string
	}{
		{s: "foobar", n: 0, want: "******"},
		{s: "foobar", n: 2, want: "fo****"},
		{s: "foobar", n: 5, want: "fooba*"},
		{s: "foobar", n: 6, want: "foobar"},
		{s: "foobar", n: 7, want: "foobar"},
	}
	for i := range cases {
		got := strutil.MaskRight(cases[i].s, cases[i].n)
		testutil.Diff(cases[i].want, got, t)
	}
}

func TestMaskLeft(t *testing.T) {
	cases := []struct {
		s    string
		n    int
		want string
	}{
		{s: "foobar", n: 0, want: "foobar"},
		{s: "foobar", n: 2, want: "****ar"},
		{s: "foobar", n: 5, want: "*oobar"},
		{s: "foobar", n: 6, want: "foobar"},
		{s: "foobar", n: 7, want: "foobar"},
	}
	for i := range cases {
		got := strutil.MaskLeft(cases[i].s, cases[i].n)
		testutil.Diff(cases[i].want, got, t)
	}
}

func TestDecapitalize(t *testing.T) {
	testutil.Diff("", strutil.Decapitalize(""), t)
	testutil.Diff("foobar", strutil.Decapitalize("Foobar"), t)
	testutil.Diff("fooBar", strutil.Decapitalize("FooBar"), t)
}

func TestIsUUID(t *testing.T) {
	testutil.Diff(false, strutil.IsUUID(""), t)
	testutil.Diff(false, strutil.IsUUID("foobar"), t)
	testutil.Diff(true, strutil.IsUUID("00000000-0000-0000-0000-000000000000"), t)
	testutil.Diff(true, strutil.IsUUID("ed811898-264d-4196-a956-0767316ff687"), t)
}
