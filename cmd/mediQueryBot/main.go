package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	handlers2 "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"med-chat-bot/cmd/mediQueryBot/internal/handlers"
	"med-chat-bot/cmd/mediQueryBot/internal/providers"
	"med-chat-bot/cmd/mediQueryBot/pkg/cfg"
	"med-chat-bot/cmd/mediQueryBot/pkg/db"
	"net/http"
)

func newMySQLConnection() *db.DB {
	_db, err := db.Connect(&db.Config{
		Driver: db.DriverMySQL,
		//Username: viper.GetString(cfg.ConfigKeyDBMySQLUsername),
		//Password: viper.GetString(cfg.ConfigKeyDBMySQLPassword),
		//Host:     viper.GetString(cfg.ConfigKeyDBMySQLHost),
		//Port:     viper.GetInt64(cfg.ConfigKeyDBMySQLPort),
		Username: "videdent_tele",
		Password: "Muaxuan2024",
		Host:     "173.252.167.20",
		Port:     3306,
		Database: viper.GetString(cfg.ConfigKeyDBMySQLDatabase)})
	if err != nil {
		log.Fatalf("Connecting to MySQL DB: %v", err)
	}
	return _db
}

func main() {
	providers.BuildContainer()
	c := providers.GetContainer()
	if c == nil {
		log.Fatalf("Container hasn't been initialized yet")
	}

	// Load environment variable
	envFile, _ := godotenv.Read(".env")
	botToken := envFile["MEDI_QUERY_BOT"]

	newMySQLConnection()

	// Create bot
	b, err := gotgbot.NewBot(botToken, &gotgbot.BotOpts{
		Client: http.Client{},
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
