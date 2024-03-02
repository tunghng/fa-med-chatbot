package main

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"log"
	"med-chat-bot/db"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/providers"
	"med-chat-bot/services/searchService"
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

type SearchHandler struct {
	searchService searchService.ISearchService
}

type SearchHandlerParams struct {
	dig.In
	SearchService searchService.ISearchService
}

func NewSearchHandler(params SearchHandlerParams) *SearchHandler {
	return &SearchHandler{
		searchService: params.SearchService,
	}
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
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewCommand("start", startCommand))
	dispatcher.AddHandler(handlers.NewCommand("help", helpCommand))
	dispatcher.AddHandler(handlers.NewCommand("about", aboutCommand))
	//dispatcher.AddHandler(handlers.NewCommand("query", queryCommand))
	dispatcher.AddHandler(handlers.NewMessage(message.All, handleMessage))

	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		log.Fatalf("Failed to start polling: %s", err)
	}
	log.Println("Bot started!")

	updater.Idle()
}

func startCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	welcomeMessage := fmt.Sprintf("👋 Welcome to MediQuery Bot! Your assistant for medical queries and feedback.\n\n" +
		"💡 I can help you find answers to medical questions, provide links to reputable sources, and collect feedback on medical information gaps. My goal is to make reliable medical information more accessible and to continuously improve based on your feedback.\n\n" +
		"Here's how you can get started:\n" +
		"- Use /query followed by your question to get medical information.\n" +
		"- Use /feedback to provide feedback or report gaps in our knowledge base.\n" +
		"- If you need help or want to learn more about how I work, just type /help.\n" +
		"- Curious about who I am and my mission? Type /about for more information on my background and how I operate.\n\n" +
		"What would you like to know today? Feel free to ask me anything!")
	_, err := ctx.EffectiveMessage.Reply(b, welcomeMessage, nil)
	return err
}

func helpCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	helpMessage := "🆘 Need some help? Here's what I can do for you:\n\n" +
		"/query [your question] - Use this command followed by your medical question to get information.\n" +
		"/feedback - Provide feedback or suggest improvements.\n" +
		"/about - Learn more about MediQuery Bot and our mission.\n" +
		"/help - Display this help message again.\n\n" +
		"Just type your command and follow the instructions. I'm here to help!"
	_, err := ctx.EffectiveMessage.Reply(b, helpMessage, nil)
	return err
}

func aboutCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	aboutMessage := "🤖 About MediQuery Bot:\n\n" +
		"MediQuery Bot is designed to provide quick, reliable medical information and collect feedback to improve our knowledge base. Our goal is to make medical information more accessible and help fill in the gaps with your feedback.\n\n" +
		"Remember: The information provided is for educational purposes and should not be considered medical advice.\n\n" +
		"👩‍💻 Developed by ___. For more information or support, contact us at ___."
	_, err := ctx.EffectiveMessage.Reply(b, aboutMessage, nil)
	return err
}

func (_this *SearchHandler) queryCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	userQuery := ctx.EffectiveMessage.Text

	webResults, err := _this.searchService.PerformSearchWebsite(userQuery)

	databaseResult, err := _this.searchService.PerformSearchWordPress(userQuery)

	replyText := "Here are your search webResults:\n"
	var count int
	for i, item := range webResults.Items {
		count = i + 1
		replyText += fmt.Sprintf("%d. [%s](%s)\n", count, item.Title, item.Link)
	}
	for i, item := range databaseResult.Items {
		if count+i+1 < 6 {
			count = i + 1
			replyText += fmt.Sprintf("%d. [%s](%s)\n", count, item.Title, item.Link)
		}
	}

	_, err = ctx.EffectiveMessage.Reply(b, replyText, &gotgbot.SendMessageOpts{ParseMode: "Markdown", DisableWebPagePreview: true})
	return err
}

func handleMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	text := "I'm sorry, I don't understand that."
	_, err := ctx.EffectiveMessage.Reply(b, text, nil)
	return err
}
