package geo

import (
	"fmt"
	"math"
	"testing"
)

const (
	epsilon = 1e-5
)

func TestHaversineDistance(t *testing.T) {
	c1 := Coord{Lat: 40.00000, Lon: 1.000000}
	c2 := Coord{Lat: 39.00000, Lon: 2.000000}
	distance := HaversineDistance(c1, c2)

	if math.Abs(distance-140.447268) > epsilon {
		fmt.Println(distance - 140.447268)
		t.Fatalf("Expected 140.447268 km, got %f", distance)
	}
}
