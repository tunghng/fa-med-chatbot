package telegram

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/services/telegram"
	"med-chat-bot/internal/handlers"
)

// FileHandler handles all requests of Refund module.
type ChatBotHandler struct {
	handlers.BaseHandler
	teleService telegram.ITelegramService
}

type chatBotHandlerParams struct {
	dig.In
	BaseHandler handlers.BaseHandler
	TeleService telegram.ITelegramService
}

func NewChatBotHandler(params chatBotHandlerParams) *ChatBotHandler {
	return &ChatBotHandler{
		BaseHandler: params.BaseHandler,
		teleService: params.TeleService,
	}
}

func (_this *ChatBotHandler) WebHook() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := _this.teleService.CallWebhook(c)
		_this.HandleResponse(c, res, err)
	}
}
