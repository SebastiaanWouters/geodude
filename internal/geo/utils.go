package geo

import "math"

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
