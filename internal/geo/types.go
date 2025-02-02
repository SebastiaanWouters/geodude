// internal/geo/types.go
package geo

import "strconv"

type Address struct {
	HouseNumber string
	Street      string
	City        string
	PostCode    string
	Country     string
	Lat         float64
	Lon         float64
}

type AddressRange struct {
	StartNumber int
	EndNumber   int
	Step        int // Usually 2 for even/odd numbering
	Street      string
	City        string
	PostCode    string
	Country     string
	StartLat    float64
	StartLon    float64
	EndLat      float64
	EndLon      float64
}

func (ar *AddressRange) Interpolate(houseNumber string) *Address {
	num, err := strconv.Atoi(houseNumber)
	if err != nil {
		return nil
	}

	// Check if number is within range and matches step pattern
	if num < ar.StartNumber || num > ar.EndNumber {
		return nil
	}
	if ar.Step > 1 && (num-ar.StartNumber)%ar.Step != 0 {
		return nil
	}

	// Calculate position ratio
	ratio := float64(num-ar.StartNumber) / float64(ar.EndNumber-ar.StartNumber)

	// Interpolate coordinates
	lat := ar.StartLat + (ar.EndLat-ar.StartLat)*ratio
	lon := ar.StartLon + (ar.EndLon-ar.StartLon)*ratio

	return &Address{
		HouseNumber: houseNumber,
		Street:      ar.Street,
		City:        ar.City,
		PostCode:    ar.PostCode,
		Country:     ar.Country,
		Lat:         lat,
		Lon:         lon,
	}
}

type GeoIndex struct {
	Addresses     map[string]*Address      // Key: "street:housenumber:postcode"
	AddressRanges map[string]*AddressRange // Key: "street:postcode"
	StreetIndex   *QuadTree                // For spatial queries
}
