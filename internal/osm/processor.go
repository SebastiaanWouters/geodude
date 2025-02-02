// internal/osm/processor.go
package osm

import (
	"context"
	"io"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// Processor defines the interface for processing OSM data
type Processor interface {
	ProcessNode(node *Node) error
	ProcessWay(way *Way) error
	ProcessRelation(relation *Relation) error
}

// StreamProcess processes OSM data in a streaming fashion
func StreamProcess(reader io.Reader, processor Processor, onlyRoutable bool) error {
	scanner := osmpbf.New(context.Background(), reader, 3)
	defer scanner.Close()

	for scanner.Scan() {
		obj := scanner.Object()
		switch v := obj.(type) {
		case *osm.Node:
			node := &Node{
				ID:   ID(v.ID),
				Lat:  v.Lat,
				Lon:  v.Lon,
				Tags: CreateTags(v.Tags),
			}
			if err := processor.ProcessNode(node); err != nil {
				return err
			}
		case *osm.Way:
			if !onlyRoutable || isRoutable(v) {
				way := &Way{
					ID:    ID(v.ID),
					Nodes: make([]ID, len(v.Nodes)),
					Tags:  CreateTags(v.Tags),
				}
				for i, node := range v.Nodes {
					way.Nodes[i] = ID(node.ID)
				}
				if err := processor.ProcessWay(way); err != nil {
					return err
				}
			}
		case *osm.Relation:
			if shouldProcessRelation(v) {
				//TODO: convert relation to our type
				//TODO: process relation
			}
		}
	}

	return scanner.Err()
}

func isRoutable(way *osm.Way) bool {
	return way.Tags.HasTag("highway") || way.Tags.HasTag("junction")
}

func shouldProcessRelation(relation *osm.Relation) bool {
	//TODO: add logic to determine if relation should be processed
	return false
}
