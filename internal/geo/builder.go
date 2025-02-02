// internal/geo/builder.go
package geo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sebastiaanwouters/geodude/internal/osm"
)

type GeoBuilder struct {
	index      *GeoIndex
	bounds     Bounds
	maxPoints  int
	streetTags map[string]bool
	nodes      map[osm.ID]*osm.Node
}

func NewGeoBuilder() *GeoBuilder {
	return &GeoBuilder{
		index: &GeoIndex{
			Addresses:     make(map[string]*Address),
			AddressRanges: make(map[string]*AddressRange),
			StreetIndex:   NewQuadTree(Bounds{MinLat: -90, MaxLat: 90, MinLon: -180, MaxLon: 180}, 50),
		},
		streetTags: map[string]bool{
			"highway":       true,
			"residential":   true,
			"service":       true,
			"living_street": true,
		},
		nodes: make(map[osm.ID]*osm.Node),
	}
}

func (b *GeoBuilder) ProcessNode(node *osm.Node) error {
	b.nodes[node.ID] = node

	if houseNumber := node.Tags.Get("addr:housenumber"); houseNumber != "" {
		addr := &Address{
			HouseNumber: houseNumber,
			Street:      node.Tags.Get("addr:street"),
			City:        node.Tags.Get("addr:city"),
			PostCode:    node.Tags.Get("addr:postcode"),
			Country:     node.Tags.Get("addr:country"),
			Lat:         node.Lat,
			Lon:         node.Lon,
		}

		key := makeAddressKey(addr.Street, addr.HouseNumber, addr.PostCode)
		b.index.Addresses[key] = addr

		b.index.StreetIndex.Insert(Point{
			Lat:  node.Lat,
			Lon:  node.Lon,
			Data: addr,
		})
	}
	return nil
}

func (b *GeoBuilder) ProcessWay(way *osm.Way) error {
	if interpolationType := way.Tags.Get("addr:interpolation"); interpolationType != "" {
		return b.processInterpolation(way, interpolationType)
	}

	if b.isStreet(way.Tags) {
		if street := way.Tags.Get("name"); street != "" {
			centerLat, centerLon := b.calculateWayCentroid(way)
			b.index.StreetIndex.Insert(Point{
				Lat:  centerLat,
				Lon:  centerLon,
				Data: street,
			})
		}
	}

	return nil
}

func (b *GeoBuilder) processInterpolation(way *osm.Way, interpolationType string) error {
	if len(way.Nodes) < 2 {
		return nil
	}

	startNode, exists := b.nodes[way.Nodes[0]]
	if !exists {
		return fmt.Errorf("start node not found")
	}

	endNode, exists := b.nodes[way.Nodes[len(way.Nodes)-1]]
	if !exists {
		return fmt.Errorf("end node not found")
	}

	startNum, err := strconv.Atoi(startNode.Tags.Get("addr:housenumber"))
	if err != nil {
		return fmt.Errorf("invalid start house number: %w", err)
	}

	endNum, err := strconv.Atoi(endNode.Tags.Get("addr:housenumber"))
	if err != nil {
		return fmt.Errorf("invalid end house number: %w", err)
	}

	step := 1
	switch interpolationType {
	case "even":
		step = 2
	case "odd":
		step = 2
	case "all":
		step = 1
	default:
		if s, err := strconv.Atoi(interpolationType); err == nil {
			step = s
		}
	}

	addressRange := &AddressRange{
		StartNumber: startNum,
		EndNumber:   endNum,
		Step:        step,
		Street:      way.Tags.Get("addr:street"),
		City:        way.Tags.Get("addr:city"),
		PostCode:    way.Tags.Get("addr:postcode"),
		Country:     way.Tags.Get("addr:country"),
		StartLat:    startNode.Lat,
		StartLon:    startNode.Lon,
		EndLat:      endNode.Lat,
		EndLon:      endNode.Lon,
	}

	key := makeStreetKey(addressRange.Street, addressRange.PostCode)
	b.index.AddressRanges[key] = addressRange

	for num := startNum; num <= endNum; num += step {
		ratio := float64(num-startNum) / float64(endNum-startNum)
		lat := startNode.Lat + (endNode.Lat-startNode.Lat)*ratio
		lon := startNode.Lon + (endNode.Lon-startNode.Lon)*ratio

		addr := &Address{
			HouseNumber: strconv.Itoa(num),
			Street:      addressRange.Street,
			City:        addressRange.City,
			PostCode:    addressRange.PostCode,
			Country:     addressRange.Country,
			Lat:         lat,
			Lon:         lon,
		}

		key := makeAddressKey(addr.Street, addr.HouseNumber, addr.PostCode)
		b.index.Addresses[key] = addr
	}

	return nil
}

func (b *GeoBuilder) ProcessRelation(relation *osm.Relation) error {
	return nil
}

func (b *GeoBuilder) GetIndex() *GeoIndex {
	return b.index
}

func (b *GeoBuilder) calculateWayCentroid(way *osm.Way) (float64, float64) {
	var sumLat, sumLon float64
	var count int

	for _, nodeID := range way.Nodes {
		if node, exists := b.nodes[nodeID]; exists {
			sumLat += node.Lat
			sumLon += node.Lon
			count++
		}
	}

	if count == 0 {
		return 0, 0
	}

	return sumLat / float64(count), sumLon / float64(count)
}

func (b *GeoBuilder) ClearNodeCache() {
	b.nodes = make(map[osm.ID]*osm.Node)
}

func makeAddressKey(street, houseNumber, postcode string) string {
	return strings.ToLower(street + ":" + houseNumber + ":" + postcode)
}

func makeStreetKey(street, postcode string) string {
	return strings.ToLower(street + ":" + postcode)
}

func (b *GeoBuilder) isStreet(tags osm.Tags) bool {
	for tag := range b.streetTags {
		if tags.Has(tag) {
			return true
		}
	}
	return false
}
