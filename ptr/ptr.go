// Package ptr implements pointer creation helpers
package ptr

// Of returns pointer to given value
//
// Deprecated: use built-in `new()` function instead (in go 1.26+)
func Of[Value any](v Value) *Value {
	return &v
}
