// Package numutil implements number helpers
package numutil

import (
	"golang.org/x/exp/constraints"
)

// Number is a constraint that covers any integer or float type
type Number interface {
	constraints.Integer | constraints.Float
}

// PercentOf calculates what percent <part> of <total>
func PercentOf[T Number](part T, total T) float64 {
	return 100 * float64(part) / float64(total)
}

// PercentDiff calculates percentage difference between two numbers
func PercentDiff[T Number](before T, after T) float64 {
	return PercentOf(after-before, before)
}

// ChangeByPercent calculates change of given number by given percent
func ChangeByPercent[T Number](value T, percent float64) float64 {
	return float64(value) + float64(value)*percent/100
}
