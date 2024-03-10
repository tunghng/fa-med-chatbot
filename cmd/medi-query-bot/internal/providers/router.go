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

		chatbotURL := v1.Group("/telegram")
		{
			chatbotURL.POST("/webhook", hs.ChatbotHandler.TestWebHook())
		}

	}
}
