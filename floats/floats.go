package floats

import (
	"math"
)

// NearlyEqual check if the given values are nearly equals.
func NearlyEqual(a float64, b float64, threshold float64) bool {
	if a == b {
		return true
	}
	return math.Abs(a-b) <= threshold
}
