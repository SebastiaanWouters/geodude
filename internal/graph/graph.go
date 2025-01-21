package graph

import (
	"github.com/sebastiaanwouters/geodude/internal/geo"
	"github.com/sebastiaanwouters/geodude/internal/osm"
)

// Node represents a node in the graph, which corresponds to an OSM node.
type Node struct {
	ID  osm.ID
	Lat float64
	Lon float64
}

// Edge represents an edge in the graph, which corresponds to an OSM way.
type Edge struct {
	From   osm.ID
	To     osm.ID
	Weight float64 // Weight can represent distance, travel time, etc.
}

// Graph represents the graph structure with nodes and edges.
type Graph struct {
	Nodes map[osm.ID]Node
	Edges map[osm.ID][]Edge
}

// NewGraph initializes a new Graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[osm.ID]Node),
		Edges: make(map[osm.ID][]Edge),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(node Node) {
	g.Nodes[node.ID] = node
}

// AddEdge adds an edge to the graph.
func (g *Graph) AddEdge(from, to osm.ID, weight float64) {
	g.Edges[from] = append(g.Edges[from], Edge{From: from, To: to, Weight: weight})
}

// ConstructGraphFromOSMData constructs a graph from the given OSMData.
func ConstructGraphFromOSMData(data *osm.OSMData) *Graph {
	graph := NewGraph()

	// Add all nodes to the graph
	for _, osmNode := range data.Nodes {
		node := Node{
			ID:  osmNode.ID,
			Lat: osmNode.Lat,
			Lon: osmNode.Lon,
		}
		graph.AddNode(node)
	}

	// Add edges based on ways
	for _, way := range data.Ways {
		for i := 0; i < len(way.Nodes)-1; i++ {
			from := way.Nodes[i]
			to := way.Nodes[i+1]

			fromNode := graph.Nodes[from]
			toNode := graph.Nodes[to]

			// Calculate the distance between the two nodes using Haversine formula
			distance := geo.HaversineDistance(geo.Coord{Lat: fromNode.Lat, Lon: fromNode.Lon}, geo.Coord{Lat: toNode.Lat, Lon: toNode.Lon})

			// Add bidirectional edges (assuming the graph is undirected)
			graph.AddEdge(from, to, distance)
			graph.AddEdge(to, from, distance)
		}
	}

	return graph
}
