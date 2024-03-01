package search

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
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

type Post struct {
	Title string `gorm:"column:post_title"`
	Link  string `gorm:"column:guid"` // Assuming 'guid' can be used as a direct link
}

func PerformSearchWebsite(query string) (*SearchResult, error) {
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

func PerformSearchWordPress(db *gorm.DB, query string) (*SearchResult, error) {
	var posts []struct {
		Title string `gorm:"column:post_title"`
		Link  string `gorm:"column:guid"` // Assuming 'guid' can be used as a direct link
	}
	err := db.Table("wplw_posts").Select("post_title as title", "guid as link").
		Where("post_title LIKE ?", "%"+query+"%").
		Where("post_status = ?", "publish").
		Limit(5).
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	// Map directly to SearchResult format
	var searchResult SearchResult
	for _, post := range posts {
		item := struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		}{
			Title: post.Title,
			Link:  post.Link,
		}
		searchResult.Items = append(searchResult.Items, item)
	}

	return &searchResult, nil
}
