package osm

import (
	"testing"
)

func TestParsePBF(t *testing.T) {
	data, err := ParsePBF("./../../data/andorra-latest.osm.pbf")
	if err != nil {
		t.Fatal(err)
	}

	// Test nodes
	if len(data.Nodes) != 471006 {
		t.Fatalf("Expected 471006 nodes, got %d", len(data.Nodes))
	}

	if data.Nodes[3854300084].Lon != 1.5963506 {
		t.Fatalf("Expected Longitude 1.5963506, got %f", data.Nodes[3854300084].Lon)
	}

	if data.Nodes[3854300084].Lat != 42.491302600000004 {
		t.Fatalf("Expected Latitude 42.491302600000004, got %f", data.Nodes[3854300084].Lat)
	}

	// Test ways
	if len(data.Ways) != 24607 {
		t.Fatalf("Expected 24607 ways, got %d", len(data.Ways))
	}

	// Check a specific way
	wayID := ID(5203906) // Replace with a known way ID from the PBF file
	way, exists := findWayByID(data.Ways, wayID)
	if !exists {
		t.Fatalf("Expected way with ID %d to exist", wayID)
	}

	// Verify the way's node IDs
	expectedNodeIDs := []ID{3655224911, 3655224917, 3655224916, 3655224915} // Replace with known node IDs for this way
	if len(way.Nodes) != len(expectedNodeIDs) {
		t.Fatalf("Expected way to have %d nodes, got %d", len(expectedNodeIDs), len(way.Nodes))
	}

	for i, nodeID := range expectedNodeIDs {
		if way.Nodes[i] != nodeID {
			t.Fatalf("Expected node ID %d at position %d, got %d", nodeID, i, way.Nodes[i])
		}
	}

	// Test relations
	if len(data.Relations) != 614 {
		t.Fatalf("Expected 614 relations, got %d", len(data.Relations))
	}

	// Check a specific relation
	relationID := ID(18) // Replace with a known relation ID from the PBF file
	relation, exists := findRelationByID(data.Relations, relationID)
	if !exists {
		t.Fatalf("Expected relation with ID %d to exist", relationID)
	}

	// Verify the relation's members
	expectedMembers := []Member{
		{
			Type: "node",
			Ref:  53376950, // Replace with known member IDs and roles
			Role: "start",
		},
		{
			Type: "way",
			Ref:  32839690, // Replace with known member IDs and roles
			Role: "both",
		},
	}

	for i, member := range expectedMembers {
		if relation.Members[i].Type != member.Type {
			t.Fatalf("Expected member type %s at position %d, got %s", member.Type, i, relation.Members[i].Type)
		}
		if relation.Members[i].Ref != member.Ref {
			t.Fatalf("Expected member ref %d at position %d, got %d", member.Ref, i, relation.Members[i].Ref)
		}
		if relation.Members[i].Role != member.Role {
			t.Fatalf("Expected member role %s at position %d, got %s", member.Role, i, relation.Members[i].Role)
		}
	}

	// Verify the relation's tags (if applicable)
	expectedRelationTags := Tags{
		{Key: "name", Value: "Section Catalonia GNR02"},
		{Key: "name:ca", Value: "Secció Catalunya GNR02"},
		{Key: "operator", Value: "www.transeurotrail.org"},
		{Key: "ref", Value: "TET:EU:ES:GNR:02:Catalonia"},
		{Key: "type", Value: "route"},
	}
	if !compareTags(relation.Tags, expectedRelationTags) {
		t.Fatalf("Expected tags %v, got %v", expectedRelationTags, relation.Tags)
	}
}

// Helper function to find a way by ID
func findWayByID(ways []Way, id ID) (Way, bool) {
	for _, way := range ways {
		if way.ID == id {
			return way, true
		}
	}
	return Way{}, false
}

// Helper function to find a relation by ID
func findRelationByID(relations []Relation, id ID) (Relation, bool) {
	for _, relation := range relations {
		if relation.ID == id {
			return relation, true
		}
	}
	return Relation{}, false
}

// Helper function to compare two Tags slices
func compareTags(tags1, tags2 Tags) bool {
	if len(tags1) != len(tags2) {
		return false
	}

	tagMap1 := make(map[string]string)
	for _, tag := range tags1 {
		tagMap1[tag.Key] = tag.Value
	}

	tagMap2 := make(map[string]string)
	for _, tag := range tags2 {
		tagMap2[tag.Key] = tag.Value
	}

	for key, value1 := range tagMap1 {
		value2, exists := tagMap2[key]
		if !exists || value1 != value2 {
			return false
		}
	}

	return true
}
