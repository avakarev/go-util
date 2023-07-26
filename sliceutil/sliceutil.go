// Package sliceutil implements slice helpers
package sliceutil

// Contains checks whenther given element is a member of given slice
func Contains[T string | int | int64 | float64](slice []T, elem T) bool {
	for _, e := range slice {
		if elem == e {
			return true
		}
	}

	return false
}
