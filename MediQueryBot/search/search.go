package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	apiKey         = "AIzaSyBzWWzys0LEqFgCPcwC5fWhkx_AQFP1KDM"
	searchEngineID = "c26f5365e4f214268"
)

type SearchResult struct {
	Items []struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	} `json:"items"`
}

func PerformSearch(query string) (*SearchResult, error) {
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=5",
		apiKey, searchEngineID, url.QueryEscape(query))

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to perform the search: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the search response: %w", err)
	}

	var results SearchResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &results, nil
}
