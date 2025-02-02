// internal/geo/geocoder_integration_test.go
package geo

import (
	"math"
	"testing"

	"github.com/sebastiaanwouters/geodude/internal/osm"
)

func TestGeocoder_Integration(t *testing.T) {
	// Create a test dataset using the GeoBuilder
	builder := NewGeoBuilder()

	// Add some test nodes
	nodes := []*osm.Node{
		{
			ID:  1,
			Lat: 42.0,
			Lon: -71.0,
			Tags: osm.Tags{
				{Key: "addr:housenumber", Value: "10"},
				{Key: "addr:street", Value: "Main Street"},
				{Key: "addr:postcode", Value: "12345"},
				{Key: "addr:city", Value: "Test City"},
				{Key: "addr:country", Value: "Test Country"},
			},
		},
		{
			ID:  2,
			Lat: 42.1,
			Lon: -71.1,
			Tags: osm.Tags{
				{Key: "addr:housenumber", Value: "15"},
				{Key: "addr:street", Value: "Oak Avenue"},
				{Key: "addr:postcode", Value: "54321"},
				{Key: "addr:city", Value: "Other City"},
				{Key: "addr:country", Value: "Test Country"},
			},
		},
		// Nodes for interpolation
		{
			ID:  3,
			Lat: 43.0,
			Lon: -72.0,
			Tags: osm.Tags{
				{Key: "addr:housenumber", Value: "1"},
				{Key: "addr:street", Value: "Pine Road"},
				{Key: "addr:postcode", Value: "98765"},
				{Key: "addr:city", Value: "Range City"},
				{Key: "addr:country", Value: "Test Country"},
			},
		},
		{
			ID:  4,
			Lat: 43.1,
			Lon: -72.1,
			Tags: osm.Tags{
				{Key: "addr:housenumber", Value: "9"},
				{Key: "addr:street", Value: "Pine Road"},
				{Key: "addr:postcode", Value: "98765"},
				{Key: "addr:city", Value: "Range City"},
				{Key: "addr:country", Value: "Test Country"},
			},
		},
	}

	for _, node := range nodes {
		err := builder.ProcessNode(node)
		if err != nil {
			t.Fatalf("Failed to process node: %v", err)
		}
	}

	// Add a way for interpolation
	way := &osm.Way{
		ID:    1,
		Nodes: []osm.ID{3, 4}, // Reference the interpolation nodes
		Tags: osm.Tags{
			{Key: "addr:interpolation", Value: "odd"},
			{Key: "addr:street", Value: "Pine Road"},
			{Key: "addr:postcode", Value: "98765"},
			{Key: "addr:city", Value: "Range City"},
			{Key: "addr:country", Value: "Test Country"},
		},
	}

	err := builder.ProcessWay(way)
	if err != nil {
		t.Fatalf("Failed to process way: %v", err)
	}

	// Get the built index
	idx := builder.GetIndex()

	// Test cases
	t.Run("Geocoding", func(t *testing.T) {
		tests := []struct {
			name        string
			street      string
			houseNumber string
			postcode    string
			wantAddr    bool
			wantLat     float64
			wantLon     float64
		}{
			{
				name:        "Exact match",
				street:      "Main Street",
				houseNumber: "10",
				postcode:    "12345",
				wantAddr:    true,
				wantLat:     42.0,
				wantLon:     -71.0,
			},
			{
				name:        "Interpolated address",
				street:      "Pine Road",
				houseNumber: "5",
				postcode:    "98765",
				wantAddr:    true,
				wantLat:     43.05,
				wantLon:     -72.05,
			},
			{
				name:        "Fuzzy match",
				street:      "Main St",
				houseNumber: "10",
				postcode:    "12345",
				wantAddr:    true,
				wantLat:     42.0,
				wantLon:     -71.0,
			},
			{
				name:        "Non-existent address",
				street:      "Fake Street",
				houseNumber: "999",
				postcode:    "00000",
				wantAddr:    false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := idx.Geocode(tt.street, tt.houseNumber, tt.postcode)
				if tt.wantAddr {
					if err != nil {
						t.Errorf("Geocode() error = %v, wantErr = false", err)
						return
					}
					if result == nil {
						t.Error("Geocode() returned nil result when address was expected")
						return
					}

					const epsilon = 0.001
					if !almostEqual(result.Lat, tt.wantLat, epsilon) || !almostEqual(result.Lon, tt.wantLon, epsilon) {
						t.Errorf("Geocode() coordinates = (%v, %v), want (%v, %v)",
							result.Lat, result.Lon, tt.wantLat, tt.wantLon)
					}
				} else {
					if err == nil {
						t.Error("Geocode() expected error for non-existent address")
					}
				}
			})
		}
	})

	t.Run("Reverse Geocoding", func(t *testing.T) {
		tests := []struct {
			name        string
			lat         float64
			lon         float64
			wantAddr    bool
			wantStreet  string
			wantNumber  string
			maxDistance float64
		}{
			{
				name:        "Exact location",
				lat:         42.0,
				lon:         -71.0,
				wantAddr:    true,
				wantStreet:  "Main Street",
				wantNumber:  "10",
				maxDistance: 0.001,
			},
			{
				name:        "Near location",
				lat:         42.001,
				lon:         -71.001,
				wantAddr:    true,
				wantStreet:  "Main Street",
				wantNumber:  "10",
				maxDistance: 0.2,
			},
			{
				name:        "Far location",
				lat:         0.0,
				lon:         0.0,
				wantAddr:    false,
				maxDistance: 0.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := idx.ReverseGeocode(tt.lat, tt.lon)
				if tt.wantAddr {
					if err != nil {
						t.Errorf("ReverseGeocode() error = %v, wantErr = false", err)
						return
					}
					if result == nil {
						t.Error("ReverseGeocode() returned nil result when address was expected")
						return
					}
					if result.Street != tt.wantStreet {
						t.Errorf("ReverseGeocode() street = %v, want %v", result.Street, tt.wantStreet)
					}
					if result.HouseNumber != tt.wantNumber {
						t.Errorf("ReverseGeocode() house number = %v, want %v", result.HouseNumber, tt.wantNumber)
					}
					if result.Distance > tt.maxDistance {
						t.Errorf("ReverseGeocode() distance = %v, want <= %v", result.Distance, tt.maxDistance)
					}
				} else {
					if result != nil {
						t.Error("ReverseGeocode() expected nil result for far location")
					}
				}
			})
		}
	})

	t.Run("Edge Cases", func(t *testing.T) {
		tests := []struct {
			name        string
			street      string
			houseNumber string
			postcode    string
			wantErr     bool
		}{
			{
				name:        "Empty street",
				street:      "",
				houseNumber: "10",
				postcode:    "12345",
				wantErr:     true,
			},
			{
				name:        "Empty house number",
				street:      "Main Street",
				houseNumber: "",
				postcode:    "12345",
				wantErr:     true,
			},
			{
				name:        "Invalid house number",
				street:      "Main Street",
				houseNumber: "abc",
				postcode:    "12345",
				wantErr:     true,
			},
			{
				name:        "Special characters in street",
				street:      "Main Street ###",
				houseNumber: "10",
				postcode:    "12345",
				wantErr:     false, // Should handle special characters
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := idx.Geocode(tt.street, tt.houseNumber, tt.postcode)
				if (err != nil) != tt.wantErr {
					t.Errorf("Geocode() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
	})
}

func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
