// Package httputil implements http helpers
package httputil

import "fmt"

// ExpectStatus returns error if actual status not equal to expected one
func ExpectStatus(want int, got int) error {
	if want != got {
		return fmt.Errorf("expected %d http response status code, but got %d", want, got)
	}
	return nil
}
