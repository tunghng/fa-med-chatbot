package crawler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/services/crawler"
	"med-chat-bot/internal/handlers"
)

type CrawlerHandler struct {
	handlers.BaseHandler
	crawlerService crawler.ICrawlerService
}

type CrawlerHandlerParams struct {
	dig.In
	BaseHandler    handlers.BaseHandler
	CrawlerService crawler.ICrawlerService
}

func NewCrawlerHandler(params CrawlerHandlerParams) *CrawlerHandler {
	return &CrawlerHandler{
		BaseHandler:    params.BaseHandler,
		crawlerService: params.CrawlerService,
	}
}

func (_this *CrawlerHandler) Crawl() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := _this.crawlerService.Crawl(c)
		_this.HandleResponse(c, resp, err)
	}
}
