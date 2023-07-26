package httputil_test

import (
	"testing"

	"github.com/avakarev/go-util/httputil"
	"github.com/avakarev/go-util/testutil"
)

type context struct {
	headers map[string]string
}

func (c context) GetHeader(s string) string {
	return c.headers[s]
}

func TestAuthBearer(t *testing.T) {
	cases := []struct {
		ctx  context
		want string
	}{
		{
			ctx:  context{headers: map[string]string{"Foo": "Bar"}},
			want: "",
		}, {
			ctx:  context{headers: map[string]string{"Authorization": ""}},
			want: "",
		}, {
			ctx:  context{headers: map[string]string{"Authorization": "qux"}},
			want: "",
		}, {
			ctx:  context{headers: map[string]string{"Authorization": "Bearer"}},
			want: "",
		}, {
			ctx:  context{headers: map[string]string{"Authorization": "Bearer:"}},
			want: "",
		}, {
			ctx:  context{headers: map[string]string{"Authorization": "Bearer: qux"}},
			want: "qux",
		}, {
			ctx:  context{headers: map[string]string{"Authorization": "bearer: qux"}},
			want: "qux",
		}}

	for i := range cases {
		got := httputil.AuthBearer(cases[i].ctx)
		want := cases[i].want
		testutil.Diff(want, got, t)
	}
}
