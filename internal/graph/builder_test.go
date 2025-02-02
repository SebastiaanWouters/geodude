package graph

import (
	"testing"

	"github.com/sebastiaanwouters/geodude/internal/osm"
)

func TestStreamProcessing(t *testing.T) {
	builder := NewGraphBuilder()
	err := osm.ParsePBF("../../data/andorra-latest.osm.pbf", true, builder)
	if err != nil {
		t.Fatal(err)
	}

	stats := builder.GetStatistics()
	if stats.NodesProcessed == 0 {
		t.Error("Expected to process some nodes")
	}
	if stats.WaysProcessed == 0 {
		t.Error("Expected to process some ways")
	}
}
