package medBot

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"io"
	dtoWordpress "med-chat-bot/internal/dtos/wordpress"
	"med-chat-bot/internal/errors"
	"med-chat-bot/internal/ginLogger"
	"med-chat-bot/internal/repositories/wordpress"
	"med-chat-bot/pkg/db"
	"med-chat-bot/pkg/meta"

	"net/http"
	"net/url"
)

type ISearchService interface {
	PerformSearchWordPress(c *gin.Context, query string) (*meta.BasicResponse, error)
}

type searchService struct {
	db         *db.DB
	searchRepo wordpress.IFaWordpressPostRepository
}

type searchServiceArgs struct {
	dig.In
	DB         *db.DB
	SearchRepo wordpress.IFaWordpressPostRepository
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

type Post struct {
	Title string `gorm:"column:post_title"`
	Link  string `gorm:"column:guid"` // Assuming 'guid' can be used as a direct link
}

func (_this searchService) PerformSearchWordPress(c *gin.Context, query string) (*meta.BasicResponse, error) {
	posts, err := _this.searchRepo.GetPostsByTitle(_this.db, query)
	if err != nil {
		//return nil, err
		ginLogger.Gin(c).Errorf("Failed when GetPostsByTitle to err: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	res := make([]dtoWordpress.SearchResult, 0)
	for _, post := range posts {
		var item dtoWordpress.SearchResult
		item.Title = post.PostTitle
		item.Link = "https://clinicalpub.com/?p=" + post.GUID
		res = append(res, item)
	}

	webResults := make([]dtoWordpress.SearchResult, 0)

	if left := 5 - len(res); len(res) != 5 {
		searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%d",
			apiKey, searchEngineID, url.QueryEscape(query), left)

		resp, err := http.Get(searchURL)
		if err != nil {
			//return nil, fmt.Errorf("failed to perform the search: %w", err)
			ginLogger.Gin(c).Errorf("Failed to perform the search: %v", err)
			return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ginLogger.Gin(c).Errorf("Failed to read the search response: %v", err)
			return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
		}

		if err := json.Unmarshal(body, &webResults); err != nil {
			ginLogger.Gin(c).Errorf("Failed to unmarshal search results: %v", err)
			return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
		}
	}
	for _, item := range webResults {
		res = append(res, item)
	}

	response := &meta.BasicResponse{
		Meta: meta.Meta{
			Code:    http.StatusOK,
			Message: "Ok",
		},
		Data: res,
	}

	return response, nil
}
