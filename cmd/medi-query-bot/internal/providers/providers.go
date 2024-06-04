package providers

import (
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/handlers"
	cbHandler "med-chat-bot/cmd/medi-query-bot/internal/handlers/telegram"
	chatbotService "med-chat-bot/cmd/medi-query-bot/internal/services/telegram"
	"med-chat-bot/internal/errors"
	"med-chat-bot/internal/ginServer"
	handlers2 "med-chat-bot/internal/handlers"
	tlRepositories "med-chat-bot/internal/repositories/telegram"
	"med-chat-bot/internal/repositories/wordpress"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/pkg/rediscmd"
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
		_ = container.Provide(newMySQLConnection, dig.Name("faquizDB"))
		_ = container.Provide(newMySQLUserTrackingConnection, dig.Name("trackingDB"))

		_ = container.Provide(rediscmd.NewRedisCmd)

		_ = container.Provide(handlers.NewHandlers)
		_ = container.Provide(cbHandler.NewChatBotHandler)
		_ = container.Provide(cbHandler.NewImageHandler)

		_ = container.Provide(chatbotService.NewTelegramService)
		_ = container.Provide(chatbotService.NewImageService)

		_ = container.Provide(tlRepositories.NewTelegramChabotRepository)
		_ = container.Provide(tlRepositories.NewTelegramChabotResponseRepository)
		_ = container.Provide(wordpress.NewWordpressPostRepository)

	}

	return container
}

// GetContainer returns an instance of Container.
func GetContainer() *dig.Container {
	return container
}
