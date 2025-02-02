// internal/geo/quadtree.go
package geo

import "math"

type QuadTree struct {
	Bounds    Bounds
	Points    []Point
	Children  [4]*QuadTree
	MaxPoints int
	Level     int
}

type Bounds struct {
	MinLat, MaxLat float64
	MinLon, MaxLon float64
}

type Point struct {
	Lat, Lon float64
	Data     interface{}
}

func NewQuadTree(bounds Bounds, maxPoints int) *QuadTree {
	return &QuadTree{
		Bounds:    bounds,
		MaxPoints: maxPoints,
		Level:     0,
		Points:    make([]Point, 0, maxPoints),
	}
}

// Contains checks if a point is within the bounds
func (b Bounds) Contains(p Point) bool {
	return p.Lat >= b.MinLat &&
		p.Lat <= b.MaxLat &&
		p.Lon >= b.MinLon &&
		p.Lon <= b.MaxLon
}

// Intersects checks if two bounds intersect
func (b Bounds) Intersects(other Bounds) bool {
	return !(other.MinLat > b.MaxLat ||
		other.MaxLat < b.MinLat ||
		other.MinLon > b.MaxLon ||
		other.MaxLon < b.MinLon)
}

// Center returns the center point of the bounds
func (b Bounds) Center() (lat, lon float64) {
	return (b.MinLat + b.MaxLat) / 2, (b.MinLon + b.MaxLon) / 2
}

// split divides the quadtree into four children
func (qt *QuadTree) split() {
	centerLat, centerLon := qt.Bounds.Center()

	// Create four children
	qt.Children[0] = NewQuadTree(Bounds{ // Northwest
		MinLat: centerLat,
		MaxLat: qt.Bounds.MaxLat,
		MinLon: qt.Bounds.MinLon,
		MaxLon: centerLon,
	}, qt.MaxPoints)

	qt.Children[1] = NewQuadTree(Bounds{ // Northeast
		MinLat: centerLat,
		MaxLat: qt.Bounds.MaxLat,
		MinLon: centerLon,
		MaxLon: qt.Bounds.MaxLon,
	}, qt.MaxPoints)

	qt.Children[2] = NewQuadTree(Bounds{ // Southwest
		MinLat: qt.Bounds.MinLat,
		MaxLat: centerLat,
		MinLon: qt.Bounds.MinLon,
		MaxLon: centerLon,
	}, qt.MaxPoints)

	qt.Children[3] = NewQuadTree(Bounds{ // Southeast
		MinLat: qt.Bounds.MinLat,
		MaxLat: centerLat,
		MinLon: centerLon,
		MaxLon: qt.Bounds.MaxLon,
	}, qt.MaxPoints)

	// Set children's level
	for i := 0; i < 4; i++ {
		qt.Children[i].Level = qt.Level + 1
	}

	// Redistribute existing points to children
	for _, p := range qt.Points {
		for i := 0; i < 4; i++ {
			if qt.Children[i].Bounds.Contains(p) {
				qt.Children[i].Insert(p)
				break
			}
		}
	}

	// Clear points from parent
	qt.Points = nil
}

func (qt *QuadTree) Insert(p Point) {
	if !qt.Bounds.Contains(p) {
		return
	}

	// If we have children, insert into appropriate child
	if qt.Children[0] != nil {
		for i := 0; i < 4; i++ {
			if qt.Children[i].Bounds.Contains(p) {
				qt.Children[i].Insert(p)
				return
			}
		}
		return
	}

	// If we haven't reached capacity, add the point
	if len(qt.Points) < qt.MaxPoints {
		qt.Points = append(qt.Points, p)
		return
	}

	// If we've reached capacity, split and redistribute
	qt.split()
	qt.Insert(p) // Re-insert the new point
}

func (qt *QuadTree) Query(bounds Bounds) []Point {
	var results []Point

	if !qt.Bounds.Intersects(bounds) {
		return results
	}

	// If we have points, check them
	for _, p := range qt.Points {
		if bounds.Contains(p) {
			results = append(results, p)
		}
	}

	// If we have children, query them
	if qt.Children[0] != nil {
		for i := 0; i < 4; i++ {
			results = append(results, qt.Children[i].Query(bounds)...)
		}
	}

	return results
}

// QueryRadius finds all points within a given radius of a center point
func (qt *QuadTree) QueryRadius(center Point, radiusKm float64) []Point {
	// Convert radius to rough lat/lon bounds
	// 1 degree of latitude is approximately 111 km
	latDelta := radiusKm / 111.0
	// 1 degree of longitude varies with latitude
	lonDelta := radiusKm / (111.0 * math.Cos(center.Lat*math.Pi/180.0))

	searchBounds := Bounds{
		MinLat: center.Lat - latDelta,
		MaxLat: center.Lat + latDelta,
		MinLon: center.Lon - lonDelta,
		MaxLon: center.Lon + lonDelta,
	}

	// First get all points within the square bounds
	candidates := qt.Query(searchBounds)

	// Then filter by actual distance
	var results []Point
	for _, p := range candidates {
		dist := HaversineDistance(
			Coord{Lat: center.Lat, Lon: center.Lon},
			Coord{Lat: p.Lat, Lon: p.Lon},
		)
		if dist <= radiusKm {
			results = append(results, p)
		}
	}

	return results
}

// Clear removes all points from the quadtree
func (qt *QuadTree) Clear() {
	qt.Points = make([]Point, 0, qt.MaxPoints)
	for i := 0; i < 4; i++ {
		qt.Children[i] = nil
	}
}

// Size returns the total number of points in the quadtree
func (qt *QuadTree) Size() int {
	count := len(qt.Points)
	if qt.Children[0] != nil {
		for i := 0; i < 4; i++ {
			count += qt.Children[i].Size()
		}
	}
	return count
}
