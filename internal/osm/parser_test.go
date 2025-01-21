package osm

import (
	"testing"
)

func TestParsePBF(t *testing.T) {

	data, err := ParsePBF("./../../data/andorra-latest.osm.pbf")

	if err != nil {
		t.Fatal(err)
	}

	if len(data.Nodes) != 471006 {
		t.Fatalf("Expected 471006 nodes, got %d", len(data.Nodes))
	}

	if data.Nodes[3854300084].Lon != 1.5963506 {
		t.Fatalf("Expected Longitude 1.5963506, got %f", data.Nodes[3854300084].Lon)
	}

	if data.Nodes[3854300084].Lat != 42.491302600000004 {
		t.Fatalf("Expected Latitude 42.491302600000004, got %f", data.Nodes[3854300084].Lat)
	}

	if len(data.Ways) != 24607 {
		t.Fatalf("Expected 24607 ways, got %d", len(data.Ways))
	}

	if len(data.Relations) != 614 {
		t.Fatalf("Expected 614 relations, got %d", len(data.Relations))
	}

}
