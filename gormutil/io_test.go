package gormutil_test

import (
	"testing"

	"github.com/avakarev/go-util/testutil"

	"github.com/avakarev/go-util/gormutil"
)

func TestFilterTables(t *testing.T) {
	cases := []struct {
		tables         []string
		include        []string
		exclude        []string
		filteredTables []string
	}{
		{
			tables:         []string{"foo", "bar", "baz", "qux", "fred", "thud"},
			filteredTables: []string{"foo", "bar", "baz", "qux", "fred", "thud"},
		}, {
			tables:         []string{"foo", "bar", "baz", "qux", "fred", "thud"},
			include:        []string{"foo", "bar"},
			filteredTables: []string{"foo", "bar"},
		}, {
			tables:         []string{"foo", "bar", "baz", "qux", "fred", "thud"},
			exclude:        []string{"foo", "bar"},
			filteredTables: []string{"baz", "qux", "fred", "thud"},
		}, {
			tables:         []string{"foo", "bar", "baz", "qux", "fred", "thud"},
			include:        []string{"foo", "bar", "baz"},
			exclude:        []string{"baz"},
			filteredTables: []string{"foo", "bar"},
		}}

	for _, tt := range cases {
		filter := gormutil.TableFilter{IncludeTables: tt.include, ExcludeTables: tt.exclude}
		testutil.Diff(tt.filteredTables, gormutil.FilterTables(tt.tables, &filter), t)
	}
}
