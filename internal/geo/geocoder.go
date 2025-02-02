// internal/geo/geocoder.go
package geo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type GeocodeResult struct {
	Address
	Distance float64
}

func (idx *GeoIndex) Geocode(street, houseNumber, postcode string) (*GeocodeResult, error) {
	// Try exact match first
	key := makeAddressKey(street, houseNumber, postcode)
	if addr, exists := idx.Addresses[key]; exists {
		return &GeocodeResult{
			Address:  *addr,
			Distance: 0,
		}, nil
	}

	// Try interpolation
	streetKey := makeStreetKey(street, postcode)
	if range_, exists := idx.AddressRanges[streetKey]; exists {
		if addr := range_.Interpolate(houseNumber); addr != nil {
			return &GeocodeResult{
				Address:  *addr,
				Distance: 0,
			}, nil
		}
	}

	// Fuzzy search
	return idx.fuzzySearch(street, houseNumber, postcode)
}

func (idx *GeoIndex) ReverseGeocode(lat, lon float64) (*GeocodeResult, error) {
	bounds := Bounds{
		MinLat: lat - 0.001,
		MaxLat: lat + 0.001,
		MinLon: lon - 0.001,
		MaxLon: lon + 0.001,
	}

	points := idx.StreetIndex.Query(bounds)
	if len(points) == 0 {
		return nil, nil
	}

	// Find closest point
	var closest *GeocodeResult
	minDist := 1000.0 // Initialize with large value

	for _, p := range points {
		if addr, ok := p.Data.(*Address); ok {
			dist := HaversineDistance(
				Coord{Lat: lat, Lon: lon},
				Coord{Lat: addr.Lat, Lon: addr.Lon},
			)
			if dist < minDist {
				minDist = dist
				closest = &GeocodeResult{
					Address:  *addr,
					Distance: dist,
				}
			}
		}
	}

	return closest, nil
}

func (idx *GeoIndex) fuzzySearch(street, houseNumber, postcode string) (*GeocodeResult, error) {
	// Normalize input
	normalizedStreet := normalizeString(street)
	normalizedPostcode := normalizeString(postcode)

	var bestMatch *GeocodeResult
	minDistance := math.MaxFloat64
	threshold := 0.55 // Similarity threshold

	// Search through all addresses
	for key, addr := range idx.Addresses {
		// Split key back into components
		parts := strings.Split(key, ":")
		if len(parts) != 3 {
			continue
		}
		streetPart := parts[0]
		postcodePart := parts[2]

		// Calculate string similarity
		streetSimilarity := calculateSimilarity(normalizedStreet, normalizeString(streetPart))
		postcodeSimilarity := calculateSimilarity(normalizedPostcode, normalizeString(postcodePart))

		// If similarity is above threshold, consider this address
		if streetSimilarity > threshold && postcodeSimilarity > threshold {
			// Convert house numbers to integers for comparison
			targetNum, err := strconv.Atoi(houseNumber)
			if err != nil {
				continue
			}
			currentNum, err := strconv.Atoi(addr.HouseNumber)
			if err != nil {
				continue
			}

			// Calculate numeric distance between house numbers
			numericDistance := math.Abs(float64(targetNum - currentNum))

			// Combine string similarity and numeric distance into a single score
			totalDistance := (1-streetSimilarity)*10 + (1-postcodeSimilarity)*5 + numericDistance*0.1

			if totalDistance < minDistance {
				minDistance = totalDistance
				bestMatch = &GeocodeResult{
					Address:  *addr,
					Distance: totalDistance,
				}
			}
		}
	}

	if bestMatch == nil {
		return nil, fmt.Errorf("no matching address found")
	}

	return bestMatch, nil
}

// Helper functions for fuzzy search
func normalizeString(s string) string {
	// Convert to lowercase and remove common prefixes/suffixes
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)

	// Remove special characters
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, s)

	return s
}

func calculateSimilarity(s1, s2 string) float64 {
	// Use Levenshtein distance for string similarity
	distance := levenshteinDistance(s1, s2)
	maxLen := math.Max(float64(len(s1)), float64(len(s2)))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(distance)/maxLen
}

func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill in the rest of the matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = min(
					matrix[i-1][j]+1,   // deletion
					matrix[i][j-1]+1,   // insertion
					matrix[i-1][j-1]+1, // substitution
				)
			}
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(numbers ...int) int {
	result := numbers[0]
	for _, num := range numbers[1:] {
		if num < result {
			result = num
		}
	}
	return result
}
