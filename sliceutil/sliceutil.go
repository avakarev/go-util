// Package sliceutil implements slice helpers
package sliceutil

import "math/rand"

// Contains checks whenther given element is a member of given slice
//
// Deprecated: use `Contains` from builtin `slices` pkg instead (in go 1.21+)
func Contains[T string | int | int64 | float64](slice []T, elem T) bool {
	for _, e := range slice {
		if elem == e {
			return true
		}
	}

	return false
}

// Shuffle returns shuffled copy of given slice
func Shuffle[T any](slice []T) []T {
	shuffled := make([]T, len(slice))
	_ = copy(shuffled, slice)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}
