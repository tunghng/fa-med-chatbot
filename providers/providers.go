package providers

import (
	"go.uber.org/dig"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/repositories"
	"med-chat-bot/services/searchService"
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
		_ = container.Provide(searchService.NewSearchService)
		_ = container.Provide(repositories.NewLinkRepository)
	}

	return container
}

// GetContainer returns an instance of Container.
func GetContainer() *dig.Container {
	return container
}
