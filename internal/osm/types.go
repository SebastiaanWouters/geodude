package osm

import "github.com/paulmach/osm"

// NodeID is a custom type for representing OSM node IDs.
type ID int64

// Node represents an OSM node with a custom NodeID.
type Node struct {
	ID   ID
	Lat  float64
	Lon  float64
	Tags Tags
}

// Way represents an OSM way with a list of NodeIDs.
type Way struct {
	ID    ID
	Nodes []ID
	Tags  Tags
}

// Relation represents an OSM relation.
type Relation struct {
	ID      ID
	Tags    Tags
	Members []Member
}

// Member represents a member of an OSM relation.
type Member struct {
	Type string // "node", "way", or "relation"
	Ref  ID     // Reference to the member's ID
	Role string // Role of the member in the relation
}

// OSMData holds the parsed OSM data using custom types.
type OSMData struct {
	Nodes     map[ID]Node
	Ways      []Way
	Relations []Relation
}

type Tag struct {
	Key   string
	Value string
}

type Tags []Tag

func (t Tags) Get(key string) string {
	for _, tag := range t {
		if tag.Key == key {
			return tag.Value
		}
	}
	return ""
}

func (t Tags) Has(key string) bool {
	for _, tag := range t {
		if tag.Key == key {
			return true
		}
	}
	return false
}

func CreateTag(key string, value string) Tag {
	return Tag{Key: key, Value: value}
}

func CreateTags(osmTags osm.Tags) Tags {
	tags := make(Tags, len(osmTags))
	for _, tag := range osmTags {
		tags = append(tags, CreateTag(tag.Key, tag.Value))
	}
	return tags
}
