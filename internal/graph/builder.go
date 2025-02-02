// internal/graph/builder.go
package graph

import (
	"github.com/sebastiaanwouters/geodude/internal/geo"
	"github.com/sebastiaanwouters/geodude/internal/osm"
)

// GraphBuilder builds a graph from OSM data using streaming
type GraphBuilder struct {
	nodes     map[osm.ID]Node
	edges     map[osm.ID][]Edge
	nodeCount int
	wayCount  int
}

func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		nodes: make(map[osm.ID]Node),
		edges: make(map[osm.ID][]Edge),
	}
}

// ProcessNode implements osm.Processor
func (b *GraphBuilder) ProcessNode(node *osm.Node) error {
	// Only store nodes that are part of ways
	b.nodeCount++
	return nil
}

// ProcessWay implements osm.Processor
func (b *GraphBuilder) ProcessWay(way *osm.Way) error {
	b.wayCount++

	// Process nodes in the way
	for i := 0; i < len(way.Nodes)-1; i++ {
		from := way.Nodes[i]
		to := way.Nodes[i+1]

		// Add nodes if they don't exist
		if _, exists := b.nodes[from]; !exists {
			b.nodes[from] = Node{ID: from}
		}
		if _, exists := b.nodes[to]; !exists {
			b.nodes[to] = Node{ID: to}
		}

		// Calculate distance and add edges
		distance := geo.HaversineDistance(
			geo.Coord{Lat: b.nodes[from].Lat, Lon: b.nodes[from].Lon},
			geo.Coord{Lat: b.nodes[to].Lat, Lon: b.nodes[to].Lon},
		)

		b.edges[from] = append(b.edges[from], Edge{From: from, To: to, Weight: distance})
		b.edges[to] = append(b.edges[to], Edge{From: to, To: from, Weight: distance})
	}
	return nil
}

// ProcessRelation implements osm.Processor
func (b *GraphBuilder) ProcessRelation(relation *osm.Relation) error {
	// Process only relations that are relevant for routing
	return nil
}

// Build returns the final graph
func (b *GraphBuilder) Build() *Graph {
	return &Graph{
		Nodes: b.nodes,
		Edges: b.edges,
	}
}

type Statistics struct {
	NodesProcessed int
	WaysProcessed  int
	NodesInGraph   int
	EdgesInGraph   int
}

func (b *GraphBuilder) GetStatistics() Statistics {
	return Statistics{
		NodesProcessed: b.nodeCount,
		WaysProcessed:  b.wayCount,
		NodesInGraph:   len(b.nodes),
		EdgesInGraph:   len(b.edges),
	}
}
