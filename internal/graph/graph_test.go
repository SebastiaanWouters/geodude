package graph

import (
	"testing"

	"github.com/sebastiaanwouters/geodude/internal/osm"
)

func TestGraph(t *testing.T) {
	// Create a sample OSMData structure for testing
	osmData := &osm.OSMData{
		Nodes: map[osm.ID]osm.Node{
			1: {ID: 1, Lat: 52.5200, Lon: 13.4050}, // Berlin
			2: {ID: 2, Lat: 48.8566, Lon: 2.3522},  // Paris
			3: {ID: 3, Lat: 51.5074, Lon: -0.1278}, // London
		},
		Ways: []osm.Way{
			{
				ID:    101,
				Nodes: []osm.ID{1, 2}, // Berlin to Paris
			},
			{
				ID:    102,
				Nodes: []osm.ID{2, 3}, // Paris to London
			},
		},
	}

	// Construct the graph from the OSM data
	graph := ConstructGraphFromOSMData(osmData)

	// Test if all nodes were added to the graph
	if len(graph.Nodes) != 3 {
		t.Fatalf("Expected 3 nodes in the graph, got %d", len(graph.Nodes))
	}

	// Test if the nodes have the correct coordinates
	berlinNode := graph.Nodes[1]
	if berlinNode.Lat != 52.5200 || berlinNode.Lon != 13.4050 {
		t.Fatalf("Expected Berlin node coordinates (52.5200, 13.4050), got (%f, %f)", berlinNode.Lat, berlinNode.Lon)
	}

	parisNode := graph.Nodes[2]
	if parisNode.Lat != 48.8566 || parisNode.Lon != 2.3522 {
		t.Fatalf("Expected Paris node coordinates (48.8566, 2.3522), got (%f, %f)", parisNode.Lat, parisNode.Lon)
	}

	londonNode := graph.Nodes[3]
	if londonNode.Lat != 51.5074 || londonNode.Lon != -0.1278 {
		t.Fatalf("Expected London node coordinates (51.5074, -0.1278), got (%f, %f)", londonNode.Lat, londonNode.Lon)
	}

	// Test if edges were correctly added
	if len(graph.Edges[1]) != 1 {
		t.Fatalf("Expected 1 edge from Berlin, got %d", len(graph.Edges[1]))
	}

	if len(graph.Edges[2]) != 2 {
		t.Fatalf("Expected 2 edges from Paris, got %d", len(graph.Edges[2]))
	}

	if len(graph.Edges[3]) != 1 {
		t.Fatalf("Expected 1 edge from London, got %d", len(graph.Edges[3]))
	}

	// Test edge weights (distances)
	berlinToParis := graph.Edges[1][0]
	expectedDistance := 878.0 // Approximate distance between Berlin and Paris in km
	if berlinToParis.Weight < expectedDistance-10 || berlinToParis.Weight > expectedDistance+10 {
		t.Fatalf("Expected distance between Berlin and Paris to be around 878 km, got %f", berlinToParis.Weight)
	}

	parisToLondon := graph.Edges[2][1]
	expectedDistance = 344.0 // Approximate distance between Paris and London in km
	if parisToLondon.Weight < expectedDistance-10 || parisToLondon.Weight > expectedDistance+10 {
		t.Fatalf("Expected distance between Paris and London to be around 344 km, got %f", parisToLondon.Weight)
	}
}
