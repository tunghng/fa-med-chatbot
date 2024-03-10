package handlers

import (
	"go.uber.org/dig"
	cb_handlers "med-chat-bot/cmd/medi-query-bot/internal/handlers/telegram"
)

// Handlers contains all handlers.
type Handlers struct {
	ChatbotHandler *cb_handlers.ChatBotHandler
}

// NewHandlersParams contains all dependencies of handlers.
type handlersParams struct {
	dig.In
	ChatbotHandler *cb_handlers.ChatBotHandler
}

// NewHandlers returns new instance of Handlers.
func NewHandlers(params handlersParams) *Handlers {
	return &Handlers{
		ChatbotHandler: params.ChatbotHandler,
	}
}
