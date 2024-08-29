package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// ParseFeed takes a URL and a schema (empty interface) as input and unmarshals the XML feed into the provided schema.
func ParseFeed(url string, schema interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch the feed: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %v", err)
	}

	err = xml.Unmarshal(body, schema)
	if err != nil {
		return fmt.Errorf("failed to unmarshal XML: %v", err)
	}

	return nil
}
