package geo

import (
	"math"
)

type Coord struct {
	Lat float64
	Lon float64
}

const (
	EarthRadius = 6371 // Earth radius in km
)

func HaversineDistance(c1, c2 Coord) float64 {
	const R = EarthRadius
	dLat := degreesToRadians(c2.Lat - c1.Lat)
	dLon := degreesToRadians(c2.Lon - c1.Lon)
	lat1 := degreesToRadians(c1.Lat)
	lat2 := degreesToRadians(c2.Lat)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
