package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"log"
	"net/http"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	envFile, _ := godotenv.Read(".env")
	botToken := envFile["MEDI_REQUEST_BOT"]

	// Create bot
	b, err := gotgbot.NewBot(botToken, &gotgbot.BotOpts{
		Client: http.Client{},
	})

	if err != nil {
		log.Fatalf("failed to create new bot: %s", err)
	}

	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	// Register command handlers
	dispatcher.AddHandler(handlers.NewCommand("start", startCommand))
	dispatcher.AddHandler(handlers.NewCommand("help", helpCommand))
	dispatcher.AddHandler(handlers.NewCommand("about", aboutCommand))
	dispatcher.AddHandler(handlers.NewCommand("query", queryCommand))
	dispatcher.AddHandler(handlers.NewMessage(message.All, handleMessage))

	// Start receiving updates
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		log.Fatalf("Failed to start polling: %s", err)
	}
	log.Println("Bot started!")

	// Idle until the program is closed
	updater.Idle()
}

func startCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	welcomeMessage := `👋 Welcome to MediQuery Bot! Your assistant for medical queries and feedback. ...`
	_, err := ctx.EffectiveMessage.Reply(b, welcomeMessage, nil)
	return err
}

func helpCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	helpMessage := `🆘 Need some help? Here's what I can do for you: ...`
	_, err := ctx.EffectiveMessage.Reply(b, helpMessage, nil)
	return err
}

func aboutCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	aboutMessage := `🤖 About MediQuery Bot: ...`
	_, err := ctx.EffectiveMessage.Reply(b, aboutMessage, nil)
	return err
}

func queryCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	// Extracting the query from the command
	userQuery := strings.Join(ctx.Args(), " ")
	dummyResponse := `🔍 You asked: '` + userQuery + `' ...`
	_, err := ctx.EffectiveMessage.Reply(b, dummyResponse, nil)
	return err
}

func handleMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	text := ctx.EffectiveMessage.Text
	response := handleResponse(text)
	_, err := ctx.EffectiveMessage.Reply(b, response, nil)
	return err
}

func handleResponse(text string) string {
	// Handle text here
	return "I'm sorry, I don't understand that."
}
