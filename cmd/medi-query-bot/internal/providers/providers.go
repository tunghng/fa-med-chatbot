package providers

import (
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/handlers"
	medBot2 "med-chat-bot/cmd/medi-query-bot/internal/handlers/medBot"
	"med-chat-bot/cmd/medi-query-bot/internal/services/medBot"
	"med-chat-bot/internal/repositories/wordpress"
	"med-chat-bot/pkg/cfg"
)

func init() {
	cfg.SetupConfig()
}

// container is a global Container.
var container *dig.Container

// BuildContainer build all necessary containers.
func BuildContainer() *dig.Container {
	container = dig.New()
	{
		_ = container.Provide(newCfgReader)
		_ = container.Provide(newMySQLConnection)
		_ = container.Provide(medBot.NewSearchService)
		_ = container.Provide(wordpress.NewWordpressPostRepository)
		_ = container.Provide(medBot2.NewSearchHandler)
		_ = container.Provide(handlers.NewHandlers)

	}

	return container
}

// GetContainer returns an instance of Container.
func GetContainer() *dig.Container {
	return container
}
