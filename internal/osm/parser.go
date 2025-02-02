package osm

import (
	"fmt"
	"io"
	"os"
	"strings"

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

func ParsePBF(filePath string, onlyRoutable bool, processor Processor) error {
	if !strings.HasSuffix(strings.ToLower(filePath), ".osm.pbf") {
		return fmt.Errorf("invalid file extension: file must end with .osm.pbf")
	}

	var reader io.ReadCloser
	var err error

	if parsedURL, isURL := isURL(filePath); isURL {
		reader, err = getURLReader(parsedURL.String())
	} else {
		reader, err = getFileReader(filePath)
	}
	if err != nil {
		return err
	}
	defer reader.Close()

	return StreamProcess(reader, processor, onlyRoutable)
}

func getURLReader(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download file: HTTP status %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func getFileReader(filePath string) (io.ReadCloser, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}
	return os.Open(filePath)
}
