package osm

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"

	"net/http"
	"net/url"
)

func isURL(filePath string) (*url.URL, bool) {
	u, err := url.Parse(filePath)
	if err != nil {
		return nil, false
	}
	if u.Scheme != "" && u.Host != "" {
		return u, true
	}
	return nil, false
}

func ParsePBF(filePath string, onlyRoutable bool) (*OSMData, error) {
	// Check if the file has the correct extension
	if !strings.HasSuffix(strings.ToLower(filePath), ".pbf") {
		return nil, fmt.Errorf("invalid file extension: file must end with .osm.pbf")
	}

	// Check if the input is a URL
	parsedURL, isURL := isURL(filePath)
	var file *os.File
	var err error

	if isURL {
		// Download the file from the URL
		resp, err := http.Get(parsedURL.String())
		if err != nil {
			return nil, fmt.Errorf("failed to download file from URL: %w", err)
		}
		defer resp.Body.Close()

		// Check response status code
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download file: HTTP status %d", resp.StatusCode)
		}

		// Check Content-Type if available
		contentType := resp.Header.Get("Content-Type")
		if contentType != "" && !strings.Contains(contentType, "application/octet-stream") {
			return nil, fmt.Errorf("invalid content type: %s", contentType)
		}

		// Create a temporary file to read the data
		tempFile, err := os.CreateTemp("", "osm-*.osm.pbf")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary file: %w", err)
		}
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name()) // Clean up temporary file
		}()

		// Copy the response body to the temporary file
		_, err = io.Copy(tempFile, resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to copy data to temporary file: %w", err)
		}

		// Open the temporary file for reading
		file, err = os.Open(tempFile.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to open temporary file: %w", err)
		}
	} else {
		// Check if the local file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", filePath)
		}

		// Open the local file
		file, err = os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
	}
	defer file.Close()

	// Read the first few bytes to verify it's a PBF file
	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}

	// Reset file pointer to beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Proceed with parsing as before
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
		return nil, fmt.Errorf("error scanning PBF file: %w", scanErr)
	}

	return data, nil
}
