package providers

import (
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/handlers"
	medBot2 "med-chat-bot/cmd/medi-query-bot/internal/handlers/medBot"
	cbHandler "med-chat-bot/cmd/medi-query-bot/internal/handlers/telegram"
	"med-chat-bot/cmd/medi-query-bot/internal/services/medBot"
	chatbotService "med-chat-bot/cmd/medi-query-bot/internal/services/telegram"
	"med-chat-bot/internal/errors"
	"med-chat-bot/internal/ginServer"
	handlers2 "med-chat-bot/internal/handlers"
	tlRepositories "med-chat-bot/internal/repositories/telegram"
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
		_ = container.Provide(newServerConfig)
		_ = container.Provide(handlers2.NewBaseHandler)
		_ = container.Provide(newErrorParserConfig)
		_ = container.Provide(newGinEngine)
		_ = container.Provide(errors.NewErrorParser)
		_ = container.Provide(ginServer.NewGinServer)

		_ = container.Provide(setupRouter)
		_ = container.Provide(newMySQLConnection)
		_ = container.Provide(medBot.NewSearchService)
		_ = container.Provide(wordpress.NewWordpressPostRepository)
		_ = container.Provide(medBot2.NewSearchHandler)
		_ = container.Provide(handlers.NewHandlers)

		_ = container.Provide(tlRepositories.NewTelegramChabotRepository)
		_ = container.Provide(chatbotService.NewTelegramService)
		_ = container.Provide(cbHandler.NewChatBotHandler)
		_ = container.Provide(newMySQLUserTrackingConnection, dig.Name("trackingDB"))
	}

	return container
}

// GetContainer returns an instance of Container.
func GetContainer() *dig.Container {
	return container
}
