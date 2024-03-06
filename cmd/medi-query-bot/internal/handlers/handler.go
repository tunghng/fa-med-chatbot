package handlers

import (
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/handlers/medBot"
)

// Handlers contains all handlers.
type Handlers struct {
	SearchHandler *medBot.SearchHandler
}

// NewHandlersParams contains all dependencies of handlers.
type handlersParams struct {
	dig.In
	SearchHanlder *medBot.SearchHandler
}

// NewHandlers returns new instance of Handlers.
func NewHandlers(params handlersParams) *Handlers {
	return &Handlers{
		SearchHandler: params.SearchHanlder,
	}
}
