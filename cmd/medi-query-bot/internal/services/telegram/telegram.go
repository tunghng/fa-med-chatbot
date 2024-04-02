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
	"strconv"

	tlRepositories "med-chat-bot/internal/repositories/telegram"
	wprepositories "med-chat-bot/internal/repositories/wordpress"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ITelegramService interface {
	LogUserFeedback(c *gin.Context) (*meta.BasicResponse, error)
	CallWebhook(c *gin.Context) (*meta.BasicResponse, error)
	SendMessage(chatID int64, message string, messageId int) (string, error)
	FeedbackCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
	ContactCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
	StartCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
	SpecialtiesCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
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
	commandToSite := map[string][]string{
		"/r":   []string{viper.GetString(cfg.GoogleSeachEngineIDR), "radiologykey", "radiologykey.com"},
		"/m":   []string{viper.GetString(cfg.GoogleSeachEngineIDM), "musculoskeletalkey", "musculoskeletalkey.com"},
		"/p":   []string{viper.GetString(cfg.GoogleSeachEngineIDP), "plasticsurgerykey", "plasticsurgerykey.com"},
		"/d":   []string{viper.GetString(cfg.GoogleSeachEngineIDD), "pocketdentistry", "pocketdentistry.com"},
		"/t":   []string{viper.GetString(cfg.GoogleSeachEngineIDT), "thoracickey", "thoracickey.com"},
		"/v":   []string{viper.GetString(cfg.GoogleSeachEngineIDV), "veteriankey", "veteriankey.com"},
		"/n":   []string{viper.GetString(cfg.GoogleSeachEngineIDN), "neupsykey", "neupsykey.com"},
		"/nu":  []string{viper.GetString(cfg.GoogleSeachEngineIDNU), "nursekey", "nursekey.com"},
		"/o":   []string{viper.GetString(cfg.GoogleSeachEngineIDO), "obgynkey", "obgynkey.com"},
		"/on":  []string{viper.GetString(cfg.GoogleSeachEngineIDON), "oncohemakey", "oncohemakey.com"},
		"/e":   []string{viper.GetString(cfg.GoogleSeachEngineIDE), "entokey", "entokey.com"},
		"/c":   []string{viper.GetString(cfg.GoogleSeachEngineIDC), "clemedicine", "clemedicine.com"},
		"/b":   []string{viper.GetString(cfg.GoogleSeachEngineIDB), "basicmedicalkey", "basicmedicalkey.com"},
		"/a":   []string{viper.GetString(cfg.GoogleSeachEngineIDA), "aneskey", "aneskey.com"},
		"/ab":  []string{viper.GetString(cfg.GoogleSeachEngineIDAB), "abdominalkey", "abdominalkey.com"},
		"/web": []string{viper.GetString(cfg.GoogleSeachEngineID), "web", ""},
		"/app": []string{viper.GetString(cfg.GoogleSeachEngineID), "clinicalpub", "https://clinicalpub.com/"},
	}

	icons := map[string]string{
		"0": "📕",
		"1": "📙",
		"2": "📗",
		"3": "📘",
		"4": "📓",
	}

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
		messageId := update.Message.MessageID
		temp := strings.Fields(userMessage)[0]
		if temp == "/start" {
			response, action, _ = _this.StartCommand(c, chatID, messageId)
		} else if temp == "/contact" {
			response, action, _ = _this.ContactCommand(c, chatID, messageId)
		} else if temp == "/fb" {
			response, action, _ = _this.FeedbackCommand(c, chatID, messageId)
		} else if temp == "/specialties" {
			response, action, _ = _this.SpecialtiesCommand(c, chatID, messageId)
		} else if value, exists := commandToSite[temp]; exists {
			_, message := TakeQueryAndAction(userMessage)
			if message == "" {
				replyText := "Your query is empty!"
				responseTele, err := _this.SendMessage(chatID, replyText, messageId)
				if err != nil {
					log.Printf("Error sending response messagev: %", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
					return nil, err
				}
				response = responseTele
			} else {
				replyText := ""
				results, _ := _this.PerformSearch(c, message, value[0], temp == "/app")
				if len(results) == 0 {
					emptyMessage := "Sorry! But we can not find any articles that matches your query! Please try another query!"
					responseTele, err := _this.SendMessage(chatID, emptyMessage, messageId)
					if err != nil {
						log.Printf("Error sending response messagev: %", err)
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
						return nil, err
					}
					response = responseTele
				} else {
					for index, item := range results {
						replyText += fmt.Sprintf("<b>%s %s</b> \n 🌱<a href=\"%s\"><i><u>%s</u></i></a>    🔗 <a href=\"%s\"><b><u>View</u></b></a> \n\n", icons[strconv.Itoa(index)], item.Title, value[2], value[1], item.Link)
					}
					log.Printf("Message from user ID %d in chat ID %d: %s\n\n", userID, chatID, update.Message.Text)
					responseTele, err := _this.SendMessage(chatID, replyText, messageId)
					if err != nil {
						log.Printf("Error sending response messagev: %", err)
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
						return nil, err
					}
					response = responseTele
				}
			}
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
		//userID := update.CallbackQuery.From.ID
		//var chatID int64
		//if update.CallbackQuery.Message != nil {
		//	chatID = update.CallbackQuery.Message.Chat.ID
		//}
		//log.Printf("Callback query from user ID %d in chat ID %d: %s\n", userID, chatID, update.CallbackQuery.Data)
		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.Chat.ID
		callbackData := update.CallbackQuery.Data

		log.Printf("Callback query from user ID %d in chat ID %d: %s\n", userID, chatID, callbackData)

		// Determine the action based on callbackData and execute it
		switch callbackData {
		case "next":
			// Handle next page logic
			_this.HandleNextPage(chatID, update.CallbackQuery.Message.MessageID)
		case "previous":
			// Handle previous page logic
			_this.HandlePreviousPage(chatID, update.CallbackQuery.Message.MessageID)
		}

		// Acknowledge the callback query to prevent the loading state on the button
		_this.AcknowledgeCallbackQuery(update.CallbackQuery.ID)
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
func (_this *telegramService) SendMessage(chatID int64, message string, messageId int) (string, error) {
	// API endpoint for sending messages
	botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	sendMessageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	inlineKeyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "Previous", "callback_data": "previous"},
				{"text": "Next", "callback_data": "next"},
			},
		},
	}
	messageBody, err := json.Marshal(map[string]interface{}{
		"chat_id":      chatID,
		"text":         message,
		"parse_mode":   "HTML",
		"reply_markup": inlineKeyboard,
		"reply_parameters": map[string]interface{}{
			"message_id":                  messageId,
			"allow_sending_without_reply": true,
		},
	})
	if err != nil {
		return "", err
	}

	response, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response from Telegram:", err)
		return "", err
	}
	log.Printf("Response from Telegram: %s\n", string(responseBody))

	return string(responseBody), nil
}

func (_this *telegramService) FeedbackCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
	action := "REPORT"
	feedbackMessage := "OMG HELP ME DON'T REPORT US!!!"
	responseTele, err := _this.SendMessage(chatID, feedbackMessage, messageId)
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

func (_this *telegramService) SpecialtiesCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
	aboutMessage := "Specialties List:\n\n" +
		"/r <keyword> - Radiology\n" +
		"/m <keyword> - Musculoskeletal\n" +
		"/p <keyword> - Plastic Surgery\n" +
		"/d <keyword> - Dentistry\n" +
		"/t <keyword> - Thoracic\n" +
		"/v <keyword> - Veterinary\n" +
		"/n <keyword> - Neupsy\n" +
		"/nu <keyword> - Nurse\n" +
		"/o <keyword> - Obstetrics and Gynaecology\n" +
		"/on <keyword> - Hematology & Oncology\n" +
		"/e <keyword> - Otolaryngology & Ophthalmology\n" +
		"/c <keyword> - Medicine (General)\n" +
		"/b <keyword> - Basic medical\n" +
		"/a <keyword> - Anesthesia\n" +
		"/ab <keyword> - Abdomen\n"
	responseTele, err := _this.SendMessage(chatID, aboutMessage, messageId)
	action := "ABOUT"
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

func (_this *telegramService) ContactCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
	aboutMessage := "For inquiries or support, feel free to reach out to us through our Telegram channel or via email:\n\n" +
		"Telegram: [t.me/videdental](t.me/videdental)\n" +
		"Email: clinicalpub.team@gmail.com\n"
	responseTele, err := _this.SendMessage(chatID, aboutMessage, messageId)
	action := "ABOUT"
	if err != nil {
		log.Printf("Error sending response messagev: %", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to send response"})
		return "", "ERROR", err
	}
	return responseTele, action, nil
}

func (_this *telegramService) StartCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
	welcomeMessage := "👋 Hi! I'm Clinical Tree Bot, you can search books or scientific articles by titles.\n\n" +
		"Just type your request like you do it in Google. \n\n" +
		"Example:\n" +
		"- Textbook of Radiology\n" +
		"Other Commands:\n" +
		"- /web - change to web mode search\n" +
		"- /app - change to app mode search\n" +
		"- /specialty - choose the specialty you want search\n" +
		"- /contact - contact us\n" +
		"What would you like to know today? Feel free to ask me anything!"
	responseTele, err := _this.SendMessage(chatID, welcomeMessage, messageId)
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

func (_this *telegramService) PerformSearch(c *gin.Context, query string, site string, isApp bool) ([]dtos.SearchResult, error) {
	res := make([]dtos.SearchResult, 0)
	apiKey := getNextApiKey()
	if isApp == true {
		posts, err := _this.wordPressPostRepo.GetPostsByTitle(_this.dbWordPress, query)
		if err != nil {
			ginLogger.Gin(c).Errorf("Failed when GetPostsByTitle to err: %v", err)
		}

		for _, post := range posts {
			var item dtos.SearchResult
			if post.PostTitle != "" {
				item.Title = post.PostTitle
				item.Link = "https://clinicalpub.com/?p=" + post.GUID
				res = append(res, item)
			}
		}
	}

	var gsr dtos.GoogleSearchResponse
	if left := 5 - len(res); len(res) != 5 {
		searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%d",
			apiKey, site, url.QueryEscape(query), left)

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
		parts := strings.Split(item.Title, "|")
		if len(parts) > 0 {
			item.Title = strings.TrimSpace(parts[0])
		}
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

// AcknowledgeCallbackQuery sends an empty response to Telegram to stop the loading spinner on the button
func (_this *telegramService) AcknowledgeCallbackQuery(callbackQueryID string) {
	botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", botToken)
	payload := map[string]string{"callback_query_id": callbackQueryID}
	payloadBytes, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
}

func (_this *telegramService) HandleNextPage(chatID int64, messageID int) {
	newText := "Next page links: [Link 6](http://...), [Link 7](http://...), ..."
	fmt.Println(newText)
	_this.UpdateMessage(chatID, messageID, newText)
}

func (_this *telegramService) HandlePreviousPage(chatID int64, messageID int) {
	newText := "Previous page links: [Link 1](http://...), [Link 2](http://...), ..."
	fmt.Println(newText)
	_this.UpdateMessage(chatID, messageID, newText)
}

// UpdateMessage edits a message's text
func (_this *telegramService) UpdateMessage(chatID int64, messageID int, newText string) {
	botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", botToken)
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       newText,
		"parse_mode": "Markdown",
		// Add inline keyboard if pagination buttons need to be kept
	}
	payloadBytes, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
}
