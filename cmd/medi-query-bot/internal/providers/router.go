package providers

import (
	"github.com/gin-gonic/gin"
	"med-chat-bot/cmd/medi-query-bot/internal/handlers"
	"med-chat-bot/internal/ginServer"
)

// setupRouter setup router.
func setupRouter(hs *handlers.Handlers) ginServer.GinRoutingFn {
	return func(router *gin.Engine) {
		v1 := router.Group("/v1")

		// Setting up the /medbot routes within v1
		medbotURL := v1.Group("/medbot")
		{
			medbotURL.GET("/query", hs.SearchHandler.GetLinksMedbot())
		}

	}
}
