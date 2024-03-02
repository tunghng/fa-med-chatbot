package searchService

import (
	"encoding/json"
	"fmt"
	"go.uber.org/dig"
	"io"
	"med-chat-bot/internal/repositories"
	"med-chat-bot/pkg/db"
	"net/http"
	"net/url"
)

type ISearchService interface {
	PerformSearchWebsite(query string) (*SearchResult, error)
	PerformSearchWordPress(query string) (*SearchResult, error)
}

type searchService struct {
	db         *db.DB
	searchRepo repositories.ILinkRepository
}

type searchServiceArgs struct {
	dig.In
	DB         *db.DB
	SearchRepo repositories.ILinkRepository
}

func NewSearchService(args searchServiceArgs) ISearchService {
	return &searchService{
		db:         args.DB,
		searchRepo: args.SearchRepo,
	}
}

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

func (_this searchService) PerformSearchWebsite(query string) (*SearchResult, error) {
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

func (_this searchService) PerformSearchWordPress(query string) (*SearchResult, error) {
	posts, err := _this.searchRepo.GetPostsByTitle(_this.db, query)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	for _, post := range posts {
		result.Items = append(result.Items, struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		}{Title: post.Title, Link: post.Link})
	}

	return &result, nil
}
