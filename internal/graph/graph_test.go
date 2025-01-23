package graph

import (
	"testing"
)

func TestGraph(t *testing.T) {
	// Create a simple test graph
	graph := NewGraph()

	// Add test nodes
	node1 := Node{ID: 1, Lat: 52.5200, Lon: 13.4050} // Berlin
	node2 := Node{ID: 2, Lat: 48.8566, Lon: 2.3522}  // Paris
	node3 := Node{ID: 3, Lat: 51.5074, Lon: -0.1278} // London

	graph.AddNode(node1)
	graph.AddNode(node2)
	graph.AddNode(node3)

	// Add test edges
	graph.AddEdge(1, 2, 878.0) // Berlin to Paris
	graph.AddEdge(2, 1, 878.0) // Paris to Berlin
	graph.AddEdge(2, 3, 344.0) // Paris to London
	graph.AddEdge(3, 2, 344.0) // London to Paris

	// Test AdjacentNodes
	t.Run("Test AdjacentNodes", func(t *testing.T) {
		// Test node 1 (Berlin)
		adjacent := graph.AdjacentNodes(1)
		if len(adjacent) != 1 {
			t.Errorf("Expected 1 adjacent node for Berlin, got %d", len(adjacent))
		}
		if weight, exists := adjacent[2]; !exists || weight != 878.0 {
			t.Errorf("Expected edge weight 878.0 between Berlin and Paris, got %f", weight)
		}

		// Test node 2 (Paris)
		adjacent = graph.AdjacentNodes(2)
		if len(adjacent) != 2 {
			t.Errorf("Expected 2 adjacent nodes for Paris, got %d", len(adjacent))
		}
		if weight, exists := adjacent[1]; !exists || weight != 878.0 {
			t.Errorf("Expected edge weight 878.0 between Paris and Berlin, got %f", weight)
		}
		if weight, exists := adjacent[3]; !exists || weight != 344.0 {
			t.Errorf("Expected edge weight 344.0 between Paris and London, got %f", weight)
		}

		// Test non-existent node
		adjacent = graph.AdjacentNodes(999)
		if len(adjacent) != 0 {
			t.Errorf("Expected 0 adjacent nodes for non-existent node, got %d", len(adjacent))
		}
	})

	// Test AdjacentNodesWithData
	t.Run("Test AdjacentNodesWithData", func(t *testing.T) {
		// Test node 2 (Paris)
		adjacent := graph.AdjacentNodesWithData(2)
		if len(adjacent) != 2 {
			t.Errorf("Expected 2 adjacent nodes for Paris, got %d", len(adjacent))
		}

		// Verify first adjacent node (Berlin)
		if adjacent[0].Node.ID != 1 {
			t.Errorf("Expected first adjacent node to be Berlin (ID 1), got %d", adjacent[0].Node.ID)
		}
		if adjacent[0].Weight != 878.0 {
			t.Errorf("Expected weight 878.0 for Berlin, got %f", adjacent[0].Weight)
		}

		// Verify second adjacent node (London)
		if adjacent[1].Node.ID != 3 {
			t.Errorf("Expected second adjacent node to be London (ID 3), got %d", adjacent[1].Node.ID)
		}
		if adjacent[1].Weight != 344.0 {
			t.Errorf("Expected weight 344.0 for London, got %f", adjacent[1].Weight)
		}

		// Test non-existent node
		adjacent = graph.AdjacentNodesWithData(999)
		if len(adjacent) != 0 {
			t.Errorf("Expected 0 adjacent nodes for non-existent node, got %d", len(adjacent))
		}
	})
}
