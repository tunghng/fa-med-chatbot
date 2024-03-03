package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/dig"
	"med-chat-bot/internal/services/searchService"
)

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

func (_this *SearchHandler) StartCommand(b *gotgbot.Bot, ctx *ext.Context) error {
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

func (_this *SearchHandler) HelpCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	helpMessage := "🆘 Need some help? Here's what I can do for you:\n\n" +
		"/query [your question] - Use this command followed by your medical question to get information.\n" +
		"/feedback - Provide feedback or suggest improvements.\n" +
		"/about - Learn more about MediQuery Bot and our mission.\n" +
		"/help - Display this help message again.\n\n" +
		"Just type your command and follow the instructions. I'm here to help!"
	_, err := ctx.EffectiveMessage.Reply(b, helpMessage, nil)
	return err
}

func (_this *SearchHandler) AboutCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	aboutMessage := "🤖 About MediQuery Bot:\n\n" +
		"MediQuery Bot is designed to provide quick, reliable medical information and collect feedback to improve our knowledge base. Our goal is to make medical information more accessible and help fill in the gaps with your feedback.\n\n" +
		"Remember: The information provided is for educational purposes and should not be considered medical advice.\n\n" +
		"👩‍💻 Developed by ___. For more information or support, contact us at ___."
	_, err := ctx.EffectiveMessage.Reply(b, aboutMessage, nil)
	return err
}

func (_this *SearchHandler) QueryCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	userQuery := ctx.EffectiveMessage.Text

	databaseResult, err := _this.searchService.PerformSearchWordPress(userQuery)
	webResults, err := _this.searchService.PerformSearchWebsite(userQuery)

	replyText := "Here are your search webResults:\n"
	var count int
	for i, item := range databaseResult.Items {
		count = i + 1
		replyText += fmt.Sprintf("%d. [%s](%s)\n", count, item.Title, item.Link)
	}
	for i, item := range webResults.Items {
		if count+i+1 < 6 {
			count = i + 1
			replyText += fmt.Sprintf("%d. [%s](%s)\n", count, item.Title, item.Link)
		}
	}

	_, err = ctx.EffectiveMessage.Reply(b, replyText, &gotgbot.SendMessageOpts{ParseMode: "Markdown", DisableWebPagePreview: true})
	return err
}

func (_this *SearchHandler) HandleMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	text := "I'm sorry, I don't understand that."
	_, err := ctx.EffectiveMessage.Reply(b, text, nil)
	return err
}
