package telegram

import (
	"bytes"
	"encoding/json"
	"med-chat-bot/internal/dtos"
	"med-chat-bot/internal/errors"
	"med-chat-bot/internal/ginLogger"
	"med-chat-bot/internal/meta"
	"med-chat-bot/internal/models"
	"med-chat-bot/pkg/db"

	tlRepositories "med-chat-bot/internal/repositories/telegram"
	//wprepositories "fa-quiz-next-gen/internal/repositories/wordpress"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"io"
	"io/ioutil"
	"log"
	wprepositories "med-chat-bot/internal/repositories/wordpress"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ITelegramService interface {
	LogUserFeedback(c *gin.Context) (*meta.BasicResponse, error)
	WebHookTest(c *gin.Context) (*meta.BasicResponse, error)
	SendMessage(chatID int64, message string) error
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

func (_this *telegramService) WebHookTest(c *gin.Context) (*meta.BasicResponse, error) {
	var update dtos.TelegramUpdate
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
			helpMessage := "🆘 Need some help? Here's what I can do for you:\n\n" +
				"/query [your question] - Use this command followed by your medical question to get information.\n" +
				"/feedback - Provide feedback or suggest improvements.\n" +
				"/about - Learn more about MediQuery Bot and our mission.\n" +
				"/help - Display this help message again.\n\n" +
				"Just type your command and follow the instructions. I'm here to help!"
			log.Printf("Message from user ID %d in chat ID %d: %s\n", userID, chatID, update.Message.Text)
			if err := _this.SendMessage(chatID, helpMessage); err != nil {
				log.Printf("Error sending response messagev: %", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
				return nil, err
			}
		} else if strings.HasPrefix(userMessage, "/about") {
			aboutMessage := "🤖 About MediQuery Bot:\n\n" +
				"MediQuery Bot is designed to provide quick, reliable medical information and collect feedback to improve our knowledge base. Our goal is to make medical information more accessible and help fill in the gaps with your feedback.\n\n" +
				"Remember: The information provided is for educational purposes and should not be considered medical advice.\n\n" +
				"👩‍💻 Developed by ___. For more information or support, contact us at ___."
			log.Printf("Message from user ID %d in chat ID %d: %s\n", userID, chatID, update.Message.Text)
			if err := _this.SendMessage(chatID, aboutMessage); err != nil {
				log.Printf("Error sending response messagev: %", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
				return nil, err
			}
		} else if strings.HasPrefix(userMessage, "/start") {
			welcomeMessage := "👋 Welcome to MediQuery Bot! Your assistant for medical queries and feedback.\n\n" +
				"💡 I can help you find answers to medical questions, provide links to reputable sources, and collect feedback on medical information gaps. My goal is to make reliable medical information more accessible and to continuously improve based on your feedback.\n\n" +
				"Here's how you can get started:\n" +
				"- Use /query followed by your question to get medical information.\n" +
				"- Use /feedback to provide feedback or report gaps in our knowledge base.\n" +
				"- If you need help or want to learn more about how I work, just type /help.\n" +
				"- Curious about who I am and my mission? Type /about for more information on my background and how I operate.\n\n" +
				"What would you like to know today? Feel free to ask me anything!"
			log.Printf("Message from user ID %d in chat ID %d: %s\n", userID, chatID, update.Message.Text)
			if err := _this.SendMessage(chatID, welcomeMessage); err != nil {
				log.Printf("Error sending response messagev: %", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
				return nil, err
			}
		} else if strings.HasPrefix(userMessage, "/query") {
			_, message := TakeQueryAndAction(userMessage)
			if message == "" {
				replyText := "Your query is empty!"
				if err := _this.SendMessage(chatID, replyText); err != nil {
					log.Printf("Error sending response messagev: %", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
					return nil, err
				}
			} else {
				replyText := "Here are your search webResults:\n"
				results, _ := _this.PerformSearch(c, message)
				for index, item := range results {
					replyText += fmt.Sprintf("%d, Title: %s \n Link to the Website: %s \n", index+1, item.Title, item.Link)
				}
				log.Printf("Message from user ID %d in chat ID %d: %s\n", userID, chatID, update.Message.Text)
				if err := _this.SendMessage(chatID, replyText); err != nil {
					log.Printf("Error sending response messagev: %", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
					return nil, err
				}
			}
		} else if strings.HasPrefix(userMessage, "/feedback") {
			action = "REPORT"
			feedbackMessage := "OMG HELP ME DON'T REPORT US!!!"
			log.Printf("Message from user ID %d in chat ID %d: %s\n", userID, chatID, update.Message.Text)
			if err := _this.SendMessage(chatID, feedbackMessage); err != nil {
				log.Printf("Error sending response messagev: %", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
				return nil, err
			}
		} else {
			replyText := "🆘 I don't understand! In general, here's what I can do for you:\n\n" +
				"/query [your question] - Use this command followed by your medical question to get information.\n" +
				"/feedback - Provide feedback or suggest improvements.\n" +
				"/about - Learn more about MediQuery Bot and our mission.\n" +
				"/help - Display this help message again.\n\n" +
				"Just type your command and follow the instructions. I'm here to help!"
			log.Printf("Message from user ID %d in chat ID %d: %s\n", userID, chatID, update.Message.Text)
			if err := _this.SendMessage(chatID, replyText); err != nil {
				log.Printf("Error sending response messagev: %", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
				return nil, err
			}
		}
		userInfo := &models.UserTrackingChatBot{
			ChatID:      update.Message.Chat.ID,
			UserID:      update.Message.From.ID,
			Action:      action,
			MessageText: update.Message.Text,
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
	response := &meta.BasicResponse{
		Meta: meta.Meta{
			Code: http.StatusOK,
		},
		Data: "Success!!!",
	}
	return response, nil
}

func (_this *telegramService) PerformSearch(c *gin.Context, query string) ([]dtos.SearchResult, error) {
	//apiKey := viper.GetString(cfg.GoogleSeachEngineAPIKey)
	//searchEngineID := viper.GetString(cfg.GoogleSeachEngineID)
	apiKey := "AIzaSyBzWWzys0LEqFgCPcwC5fWhkx_AQFP1KDM"
	searchEngineID := "c26f5365e4f214268"
	posts, err := _this.wordPressPostRepo.GetPostsByTitle(_this.dbWordPress, query)
	if err != nil {
		ginLogger.Gin(c).Errorf("Failed when GetPostsByTitle to err: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	res := make([]dtos.SearchResult, 0)
	for _, post := range posts {
		var item dtos.SearchResult
		if item.Title != "" {
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
			//return nil, fmt.Errorf("failed to perform the search: %w", err)
			ginLogger.Gin(c).Errorf("Failed to perform the search: %v", err)
			return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ginLogger.Gin(c).Errorf("Failed to read the search response: %v", err)
			return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
		}

		if err := json.Unmarshal(body, &gsr); err != nil {
			ginLogger.Gin(c).Errorf("Failed to unmarshal search results: %v", err)
			return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
		}
	}
	for _, item := range gsr.Items {
		res = append(res, item)
	}
	if len(res) == 0 {
		return res, nil
	}
	convertedRes, err := _this.ConvertLink(res)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error converting links: %v\n", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	return convertedRes, nil
}

// SendMessage sends a message to a user in Telegram
func (_this *telegramService) SendMessage(chatID int64, message string) error {
	// API endpoint for sending messages
	//botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	botToken := "6327700438:AAGFtsaQcs0dT3Uiwp9idqAitA814OosEg4"
	sendMessageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	messageBody, err := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	})
	if err != nil {
		return err
	}

	response, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response from Telegram:", err)
		return err
	}
	log.Printf("Response from Telegram: %s\n", responseBody)

	return nil
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

func (_this *telegramService) ConvertLink(items []dtos.SearchResult) ([]dtos.SearchResult, error) {
	//linklyToken := cfg.ConfigLinklyToken
	linklyToken := "xUBAGyqEnKqs4WOWEub44g=="
	//linklyWorkspaceId := cfg.ConfigLinklyWorkspaceId
	linklyWorkspaceId := "181063"

	var requests []dtos.LinkConversionRequest
	for _, item := range items {
		requests = append(requests, dtos.LinkConversionRequest{
			WorkspaceID: linklyWorkspaceId,
			Url:         item.Link,
		})
	}

	requestBody, err := json.Marshal(requests)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("https://app.linklyhq.com/api/v1/workspace/%s/links?api_key=%s", linklyWorkspaceId, linklyToken)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return nil, fmt.Errorf("expected JSON response but got: %s", contentType)
	}

	var linkResponses []dtos.LinkConversionResponse
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&linkResponses); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for i, linkResponse := range linkResponses {
		if linkResponse.Status == "success" && i < len(items) {
			items[i].Link = linkResponse.Link.FullURL
		}
	}

	return items, nil
}
