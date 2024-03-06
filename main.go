package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	handlers2 "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"med-chat-bot/cmd/internal/handlers"
	"med-chat-bot/cmd/internal/providers"
)

func main() {
	providers.BuildContainer()
	c := providers.GetContainer()
	if c == nil {
		log.Fatalf("Container hasn't been initialized yet")
	}

	// Load environment variable
	envFile, _ := godotenv.Read(".env")
	botToken := envFile["MEDI_QUERY_BOT"]

	// Create bot
	b, err := gotgbot.NewBot(botToken, &gotgbot.BotOpts{
		//Client: http.Client{},
	})

	if err != nil {
		log.Fatalf("failed to create new bot: %s", err)
	}
	var searchHandler *handlers.SearchHandler
	err = c.Invoke(func(sh *handlers.SearchHandler) {
		searchHandler = sh
	})
	if err != nil {
		log.Fatalf("Failed to get SearchHandler: %s", err)
	}
	// Command Handlers
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers2.NewCommand("start", searchHandler.StartCommand))
	dispatcher.AddHandler(handlers2.NewCommand("help", searchHandler.HelpCommand))
	dispatcher.AddHandler(handlers2.NewCommand("about", searchHandler.AboutCommand))
	dispatcher.AddHandler(handlers2.NewCommand("query", searchHandler.QueryCommand))
	dispatcher.AddHandler(handlers2.NewMessage(message.All, searchHandler.HandleMessage))

	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		log.Fatalf("Failed to start polling: %s", err)
	}
	log.Println("Bot started!")

	updater.Idle()
}
