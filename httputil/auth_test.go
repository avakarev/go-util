package httputil_test

import (
	"testing"

	"github.com/avakarev/go-util/httputil"
	"github.com/avakarev/go-util/testutil"
)

func TestAuthBearer(t *testing.T) {
	cases := []struct {
		auth   string
		bearer string
	}{
		{
			auth:   "",
			bearer: "",
		}, {
			auth:   "qux",
			bearer: "",
		}, {
			auth:   "Bearer",
			bearer: "",
		}, {
			auth:   "Bearer:",
			bearer: "",
		}, {
			auth:   "Bearer: qux",
			bearer: "qux",
		}, {
			auth:   "bearer: qux",
			bearer: "qux",
		}}

	for i := range cases {
		got := httputil.AuthBearer(cases[i].auth)
		want := cases[i].bearer
		testutil.Diff(want, got, t)
	}
}
