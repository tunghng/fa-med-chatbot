package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"med-chat-bot/db"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/providers"
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

	// Command Handlers
	updater := ext.NewUpdater(nil)
	//dispatcher := updater.Dispatcher

	//dispatcher.AddHandler(handlers.NewCommand("start", startCommand))
	//dispatcher.AddHandler(handlers.NewCommand("help", helpCommand))
	//dispatcher.AddHandler(handlers.NewCommand("about", aboutCommand))
	//dispatcher.AddHandler(handlers.NewCommand("query", queryCommand))
	//dispatcher.AddHandler(handlers.NewMessage(message.All, handleMessage))

	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		log.Fatalf("Failed to start polling: %s", err)
	}
	log.Println("Bot started!")

	updater.Idle()
}
