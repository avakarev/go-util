package testutil_test

import (
	"testing"

	"github.com/avakarev/go-util/testutil"
)

type user struct {
	Name    string
	address string
}

func TestDiff(t *testing.T) {
	u1 := &user{Name: "foo", address: "addr 1"}
	u2 := &user{Name: "foo", address: "addr 2"}
	testutil.Diff(u1, u2, t)
}
