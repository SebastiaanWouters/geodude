package osm

import (
	"context"
	"os"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// ParsePBF parses an OSM PBF file and returns the data using custom types.
func ParsePBF(filePath string, onlyRoutable bool) (*OSMData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := osmpbf.New(context.Background(), file, 3)
	defer scanner.Close()

	data := &OSMData{
		Nodes:     make(map[ID]Node),
		Ways:      make([]Way, 0),
		Relations: make([]Relation, 0),
	}

	for scanner.Scan() {
		o := scanner.Object()

		switch v := o.(type) {
		case *osm.Node:
			data.Nodes[ID(v.ID)] = Node{
				ID:  ID(v.ID),
				Lat: v.Lat,
				Lon: v.Lon,
			}
		case *osm.Way:
			if !onlyRoutable {
				way := Way{
					ID:    ID(v.ID),
					Nodes: make([]ID, len(v.Nodes)),
				}
				for i, node := range v.Nodes {
					way.Nodes[i] = ID(node.ID)
				}
				data.Ways = append(data.Ways, way)
				break
			}
			if v.Tags.HasTag("highway") || v.Tags.HasTag("junction") {
				way := Way{
					ID:    ID(v.ID),
					Nodes: make([]ID, len(v.Nodes)),
				}
				for i, node := range v.Nodes {
					way.Nodes[i] = ID(node.ID)
				}
				data.Ways = append(data.Ways, way)
			}
		case *osm.Relation:
			tags := make(Tags, len(v.Tags))
			for i, tag := range v.Tags {
				tags[i] = Tag{
					Key:   tag.Key,
					Value: tag.Value,
				}
			}
			relation := Relation{
				ID:      ID(v.ID),
				Tags:    tags,
				Members: make([]Member, len(v.Members)),
			}
			for i, member := range v.Members {
				relation.Members[i] = Member{
					Type: string(member.Type),
					Ref:  ID(member.Ref),
					Role: member.Role,
				}
			}
			data.Relations = append(data.Relations, relation)
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		panic(scanErr)
	}

	return data, nil
}
