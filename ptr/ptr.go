// Package ptr implements pointer creation helpers
package ptr

// Of returns pointer to given value
func Of[Value any](v Value) *Value {
	return &v
}
