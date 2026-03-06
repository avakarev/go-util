// Package sliceutil implements slice helpers
package sliceutil

import "math/rand"

// Shuffle returns shuffled copy of given slice
func Shuffle[T any](slice []T) []T {
	shuffled := make([]T, len(slice))
	_ = copy(shuffled, slice)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}
