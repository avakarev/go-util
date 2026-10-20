// Package numutil implements number helpers
package numutil

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Number is a constraint that covers any integer or float type
type Number interface {
	constraints.Integer | constraints.Float
}

// Percent calculates value of given <percent> from given <total>
func Percent[T Number](percent T, total T) float64 {
	return float64(percent) * float64(total) / 100
}

// PercentOf calculates what percent <part> of <total>.
// It returns 0 when total is 0 to avoid an Inf/NaN result.
func PercentOf[T Number](part T, total T) float64 {
	if total == 0 {
		return 0
	}
	return 100 * float64(part) / float64(total)
}

// PercentDiff calculates percentage difference between two numbers.
// It returns 0 when before is 0 to avoid an Inf/NaN result.
func PercentDiff[T Number](before T, after T) float64 {
	return PercentOf(after-before, before)
}

// ChangeByPercent calculates change of given number by given percent
func ChangeByPercent[T Number](value T, percent float64) float64 {
	return float64(value) + float64(value)*percent/100
}

// Round rounds given float number to given precision (decimal digits)
func Round[T float32 | float64](value T, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(value)*ratio) / ratio
}

// Drawdown calculates percentage of the peak value.
// It returns 0 when high is 0 to avoid an Inf/NaN result.
func Drawdown[T Number](high T, regular T) float64 {
	return PercentOf(high-regular, high)
}
