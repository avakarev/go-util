package httputil_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/avakarev/go-util/httputil"
	"github.com/avakarev/go-util/testutil"

	"github.com/jarcoal/httpmock"
)

type item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type testClient struct {
	base httputil.BaseClient
}

func (tc *testClient) getItem(id string) (*item, error) {
	var i item
	_, err := tc.base.GetJSON("/api/items/"+id, &i)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func newTestClient() *testClient {
	return &testClient{
		base: httputil.BaseClient{
			BaseURL: "https://example.org",
			Header: http.Header{
				"User-Agent": []string{"my-test-ua-1.0"},
			},
			Timeout: 1 * time.Second,
		},
	}
}

func TestBaseClientGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://example.org/api/items/1",
		httpmock.NewStringResponder(200, `{"id": "1", "name": "foobar"}`))

	i, err := newTestClient().getItem("1")
	testutil.MustNoErr(err, t)
	testutil.Diff(&item{ID: "1", Name: "foobar"}, i, t)

}
