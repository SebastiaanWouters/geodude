package geo

import (
	"fmt"
	"testing"
)

func TestAddressRange_Interpolate(t *testing.T) {
	ar := &AddressRange{
		StartNumber: 1,
		EndNumber:   5,
		Step:        2,
		Street:      "Test Street",
		City:        "Test City",
		PostCode:    "12345",
		Country:     "Test Country",
		StartLat:    0,
		StartLon:    0,
		EndLat:      1,
		EndLon:      1,
	}

	tests := []struct {
		houseNumber string
		want        bool
	}{
		{"1", true},
		{"3", true},
		{"5", true},
		{"2", false},
		{"4", false},
		{"6", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.houseNumber, func(t *testing.T) {
			got := ar.Interpolate(tt.houseNumber)
			if (got != nil) != tt.want {
				t.Errorf("Interpolate(%q) = %v, want %v", tt.houseNumber, got != nil, tt.want)
			}
		})
	}
}

func TestGeoIndex_FuzzySearch(t *testing.T) {
	idx := &GeoIndex{
		Addresses: map[string]*Address{
			"main street:10:12345": {
				Street:      "Main Street",
				HouseNumber: "10",
				PostCode:    "12345",
			},
		},
	}

	tests := []struct {
		street      string
		houseNumber string
		postcode    string
		wantErr     bool
	}{
		{"Main Strt", "10", "12345", false},
		{"Main Street", "12", "12345", false},
		{"Completely Wrong", "10", "12345", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%s-%s", tt.street, tt.houseNumber, tt.postcode), func(t *testing.T) {
			_, err := idx.fuzzySearch(tt.street, tt.houseNumber, tt.postcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("fuzzySearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
