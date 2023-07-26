package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// MustNoErr fail the test if `err` is not nil
func MustNoErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf("%sGot unexpected error: %v", caller(), err)
	}
}

// MustErr fails the test if `want` error differs from `got` error
func MustErr(want error, got error, t *testing.T) {
	if diff := cmp.Diff(want.Error(), got.Error()); diff != "" {
		t.Errorf("%sGot unexpected error (-want +got):\n%s", caller(), diff)
	}
}
