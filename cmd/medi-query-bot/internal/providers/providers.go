package providers

import (
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/handlers"
	"med-chat-bot/cmd/medi-query-bot/internal/services/medBot"
	"med-chat-bot/internal/repositories"
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
		_ = container.Provide(repositories.NewLinkRepository)
		_ = container.Provide(handlers.NewSearchHandler)
	}

	return container
}

// GetContainer returns an instance of Container.
func GetContainer() *dig.Container {
	return container
}
