package osm

import (
	"context"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

type OSMData struct {
	Nodes     map[osm.NodeID]osm.Node
	Ways      []osm.Way
	Relations []osm.Relation
}

func ParsePBF(filePath string) (*OSMData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := osmpbf.New(context.Background(), file, 3)
	defer scanner.Close()

	data := &OSMData{
		Nodes: make(map[osm.NodeID]osm.Node),
		Ways:  make([]osm.Way, 0),
	}

	for scanner.Scan() {
		o := scanner.Object()

		switch v := o.(type) {
		case *osm.Node:
			data.Nodes[v.ID] = *v
		case *osm.Way:
			data.Ways = append(data.Ways, *v)
		case *osm.Relation:
			data.Relations = append(data.Relations, *v)
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		panic(scanErr)
	}

	return data, nil
}
