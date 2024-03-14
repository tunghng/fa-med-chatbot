package telegram

import (
	"bytes"
	"encoding/json"
	"med-chat-bot/internal/dtos"
	"med-chat-bot/internal/ginLogger"
	"med-chat-bot/internal/meta"
	"med-chat-bot/internal/models"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/pkg/db"

	tlRepositories "med-chat-bot/internal/repositories/telegram"
	wprepositories "med-chat-bot/internal/repositories/wordpress"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ITelegramService interface {
	LogUserFeedback(c *gin.Context) (*meta.BasicResponse, error)
	CallWebhook(c *gin.Context) (*meta.BasicResponse, error)
	SendMessage(chatID int64, message string) (string, error)
	FeedbackCommand(c *gin.Context, chatID int64) (string, string, error)
	AboutCommand(c *gin.Context, chatID int64) (string, string, error)
	HelpCommand(c *gin.Context, chatID int64) (string, string, error)
	StartCommand(c *gin.Context, chatID int64) (string, string, error)
	UnknownCommand(c *gin.Context, chatID int64) (string, string, error)
}

type telegramService struct {
	dbWordPress       *db.DB
	dbTracking        *db.DB
	wordPressPostRepo wprepositories.IFaWordpressPostRepository
	userTrackingRepo  tlRepositories.ITelegramChabotRepository
}

type TelegramServiceArgs struct {
	dig.In
	DBWordPress       *db.DB `name:"faquizDB"`
	DBTracking        *db.DB `name:"trackingDB"`
	WordPressPostRepo wprepositories.IFaWordpressPostRepository
	UserTrackingRepo  tlRepositories.ITelegramChabotRepository
}

func NewTelegramService(args TelegramServiceArgs) ITelegramService {
	return &telegramService{
		dbWordPress:       args.DBWordPress,
		dbTracking:        args.DBTracking,
		wordPressPostRepo: args.WordPressPostRepo,
		userTrackingRepo:  args.UserTrackingRepo,
	}
}

func (_this *telegramService) LogUserFeedback(c *gin.Context) (*meta.BasicResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (_this *telegramService) CallWebhook(c *gin.Context) (*meta.BasicResponse, error) {
	var update dtos.TelegramUpdate
	var response string
	action := "SEARCH"
	if err := c.ShouldBindJSON(&update); err != nil {
		log.Printf("Error decoding update: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding update"})
		return nil, err
	}
	// Determine the type of update
	if update.Message != nil {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID
		userMessage := update.Message.Text
		if strings.HasPrefix(userMessage, "/help") {
			response, action, _ = _this.HelpCommand(c, chatID)
		} else if strings.HasPrefix(userMessage, "/about") {
			response, action, _ = _this.AboutCommand(c, chatID)
		} else if strings.HasPrefix(userMessage, "/start") {
			response, action, _ = _this.StartCommand(c, chatID)
		} else if strings.HasPrefix(userMessage, "/feedback") {
			response, action, _ = _this.FeedbackCommand(c, chatID)
		} else if strings.HasPrefix(userMessage, "/query") {
			_, message := TakeQueryAndAction(userMessage)
			if message == "" {
				replyText := "Your query is empty!"
				responseTele, err := _this.SendMessage(chatID, replyText)
				if err != nil {
					log.Printf("Error sending response messagev: %", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
					return nil, err
				}
				response = responseTele
			} else {
				replyText := "Here are your search results:\n"
				results, _ := _this.PerformSearch(c, message)
				if len(results) == 0 {
					emptyMessage := "Sorry! But we can not find any articles that matches your query! Please try another query!"
					responseTele, err := _this.SendMessage(chatID, emptyMessage)
					if err != nil {
						log.Printf("Error sending response messagev: %", err)
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
						return nil, err
					}
					response = responseTele
				} else {
					for index, item := range results {
						replyText += fmt.Sprintf("%d, Title: %s \n Link to the Website: %s \n", index+1, item.Title, item.Link)
					}
					log.Printf("Message from user ID %d in chat ID %d: %s\n", userID, chatID, update.Message.Text)
					responseTele, err := _this.SendMessage(chatID, replyText)
					if err != nil {
						log.Printf("Error sending response messagev: %", err)
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
						return nil, err
					}
					response = responseTele
				}
			}
		} else if strings.HasPrefix(userMessage, "/feedback") {
			response, action, _ = _this.FeedbackCommand(c, chatID)
		} else {
			return nil, nil
		}

		var msg dtos.TelegramMessage
		if err := json.Unmarshal([]byte(response), &msg); err != nil {
			fmt.Println("Error parsing JSON: ", err)
			return nil, err
		}
		name := msg.Result.Chat.FirstName + " " + msg.Result.Chat.LastName
		_, message := TakeQueryAndAction(userMessage)
		userInfo := &models.UserTrackingChatBot{
			ChatID:      update.Message.Chat.ID,
			UserID:      update.Message.From.ID,
			Action:      action,
			Username:    name,
			MessageText: message,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_, err := _this.userTrackingRepo.Create(_this.dbTracking, userInfo)
		if err != nil {
			return nil, err
		}
	} else if update.CallbackQuery != nil {
		userID := update.CallbackQuery.From.ID
		var chatID int64
		if update.CallbackQuery.Message != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		}
		log.Printf("Callback query from user ID %d in chat ID %d: %s\n", userID, chatID, update.CallbackQuery.Data)
	}

	c.Status(http.StatusOK)
	responseBE := &meta.BasicResponse{
		Meta: meta.Meta{
			Code: http.StatusOK,
		},
		Data: "Success!!!",
	}
	return responseBE, nil
}

// SendMessage sends a message to a user in Telegram
func (_this *telegramService) SendMessage(chatID int64, message string) (string, error) {
	// API endpoint for sending messages
	botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	sendMessageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	messageBody, err := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	})
	if err != nil {
		return "", err
	}

	response, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response from Telegram:", err)
		return "", err
	}
	log.Printf("Response from Telegram: %s\n", string(responseBody))

	return string(responseBody), nil
}

func (_this *telegramService) HelpCommand(c *gin.Context, chatID int64) (string, string, error) {
	helpMessage := "🆘 Need some help? Here's what I can do for you:\n\n" +
		"/query [your question] - Use this command followed by your medical question to get information.\n" +
		"/feedback - Provide feedback or suggest improvements.\n" +
		"/about - Learn more about MediQuery Bot and our mission.\n" +
		"/help - Display this help message again.\n\n" +
		"Just type your command and follow the instructions. I'm here to help!"
	responseTele, err := _this.SendMessage(chatID, helpMessage)
	action := "HELP"
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

func (_this *telegramService) FeedbackCommand(c *gin.Context, chatID int64) (string, string, error) {
	action := "REPORT"
	feedbackMessage := "OMG HELP ME DON'T REPORT US!!!"
	responseTele, err := _this.SendMessage(chatID, feedbackMessage)
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

func (_this *telegramService) UnknownCommand(c *gin.Context, chatID int64) (string, string, error) {
	action := "UNKNOWN"
	replyText := "🆘 I don't understand! In general, here's what I can do for you:\n\n" +
		"/query [your question] - Use this command followed by your medical question to get information.\n" +
		"/feedback - Provide feedback or suggest improvements.\n" +
		"/about - Learn more about MediQuery Bot and our mission.\n" +
		"/help - Display this help message again.\n\n" +
		"Just type your command and follow the instructions. I'm here to help!"
	responseTele, err := _this.SendMessage(chatID, replyText)
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

func (_this *telegramService) AboutCommand(c *gin.Context, chatID int64) (string, string, error) {
	aboutMessage := "🤖 About MediQuery Bot:\n\n" +
		"MediQuery Bot is designed to provide quick, reliable medical information and collect feedback to improve our knowledge base. Our goal is to make medical information more accessible and help fill in the gaps with your feedback.\n\n" +
		"Remember: The information provided is for educational purposes and should not be considered medical advice.\n\n" +
		"👩‍💻 Developed by ___. For more information or support, contact us at ___."
	responseTele, err := _this.SendMessage(chatID, aboutMessage)
	action := "ABOUT"
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

func (_this *telegramService) StartCommand(c *gin.Context, chatID int64) (string, string, error) {
	welcomeMessage := "👋 Welcome to MediQuery Bot! Your assistant for medical queries and feedback.\n\n" +
		"💡 I can help you find answers to medical questions, provide links to reputable sources, and collect feedback on medical information gaps. My goal is to make reliable medical information more accessible and to continuously improve based on your feedback.\n\n" +
		"Here's how you can get started:\n" +
		"- Use /query followed by your question to get medical information.\n" +
		"- Use /feedback to provide feedback or report gaps in our knowledge base.\n" +
		"- If you need help or want to learn more about how I work, just type /help.\n" +
		"- Curious about who I am and my mission? Type /about for more information on my background and how I operate.\n\n" +
		"What would you like to know today? Feel free to ask me anything!"
	responseTele, err := _this.SendMessage(chatID, welcomeMessage)
	action := "START"
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

var apiKeyIndex int

func getNextApiKey() string {
	keys := []string{
		viper.GetString(cfg.GoogleSeachEngineAPIKey1),
		viper.GetString(cfg.GoogleSeachEngineAPIKey2),
		viper.GetString(cfg.GoogleSeachEngineAPIKey3),
	}
	apiKey := keys[apiKeyIndex]
	apiKeyIndex = (apiKeyIndex + 1) % len(keys) // Cycle through the keys
	return apiKey
}

func (_this *telegramService) PerformSearch(c *gin.Context, query string) ([]dtos.SearchResult, error) {
	apiKey := getNextApiKey()
	searchEngineID := viper.GetString(cfg.GoogleSeachEngineID)
	posts, err := _this.wordPressPostRepo.GetPostsByTitle(_this.dbWordPress, query)
	if err != nil {
		ginLogger.Gin(c).Errorf("Failed when GetPostsByTitle to err: %v", err)
	}
	res := make([]dtos.SearchResult, 0)
	for _, post := range posts {
		var item dtos.SearchResult
		if post.PostTitle != "" {
			item.Title = post.PostTitle
			item.Link = "https://clinicalpub.com/?p=" + post.GUID
			res = append(res, item)
		}
	}

	var gsr dtos.GoogleSearchResponse
	if left := 5 - len(res); len(res) != 5 {
		searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%d",
			apiKey, searchEngineID, url.QueryEscape(query), left)

		resp, err := http.Get(searchURL)
		if err != nil {
			ginLogger.Gin(c).Errorf("Failed to perform the search: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ginLogger.Gin(c).Errorf("Failed to read the search response: %v", err)
		}

		if err := json.Unmarshal(body, &gsr); err != nil {
			ginLogger.Gin(c).Errorf("Failed to unmarshal search results: %v", err)
		}
	}
	for _, item := range gsr.Items {
		res = append(res, item)
	}
	if len(res) == 0 {
		return res, nil
	}
	convertedRes, err := _this.ConvertLink(c, res)
	if err != nil {
		return res, nil
	}
	return convertedRes, nil
}

func (_this *telegramService) ConvertLink(c *gin.Context, items []dtos.SearchResult) ([]dtos.SearchResult, error) {
	linklyToken := viper.GetString(cfg.ConfigLinklyToken)
	linklyWorkspaceId := viper.GetString(cfg.ConfigLinklyWorkspaceId)

	var requests []dtos.LinkConversionRequest
	for _, item := range items {
		requests = append(requests, dtos.LinkConversionRequest{
			WorkspaceID: linklyWorkspaceId,
			Url:         item.Link,
		})
	}

	requestBody, err := json.Marshal(requests)
	if err != nil {
		ginLogger.Gin(c).Errorf("failed to marshal request body: %v", err)
	}

	url := fmt.Sprintf("https://app.linklyhq.com/api/v1/workspace/%s/links?api_key=%s", linklyWorkspaceId, linklyToken)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		ginLogger.Gin(c).Errorf("failed to execute request: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		ginLogger.Gin(c).Errorf("API request failed with status code: %d", response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		ginLogger.Gin(c).Errorf("expected JSON response but got: %s", contentType)
	}

	defer response.Body.Close()
	var linkResponses []dtos.LinkConversionResponse
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&linkResponses)

	for i, linkResponse := range linkResponses {
		if linkResponse.Status == "success" && i < len(items) {
			items[i].Link = linkResponse.Link.FullURL
		}
	}

	return items, err
}

func TakeQueryAndAction(userMessage string) (action string, message string) {
	parts := strings.SplitN(userMessage, " ", 2)
	if len(parts) > 0 {
		action = parts[0]
	}
	if len(parts) > 1 {
		message = parts[1]
	}
	return action, message
}
