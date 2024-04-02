package telegram

import (
	"bytes"
	"encoding/json"
	"med-chat-bot/internal/dtos"
	"med-chat-bot/internal/errors"
	"med-chat-bot/internal/ginLogger"
	"med-chat-bot/internal/meta"
	"med-chat-bot/internal/models"
	"med-chat-bot/pkg/cfg"
	"med-chat-bot/pkg/db"
	"med-chat-bot/pkg/rediscmd"
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
	CallWebhook(c *gin.Context) (*meta.BasicResponse, error)
	SendMessage(c *gin.Context, chatID int64, message string, messageId int) (string, error)
	ContactCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
	StartCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
	ModeCommand(c *gin.Context, chatID int64, messageId int, mode string) (string, string, error)
	SpecialtiesCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
	SendMessageWithCallBack(c *gin.Context, chatID int64, message string, messageId int, currentPage int, globalUpdate dtos.InteractionContext) (string, error)
	// FeedbackCommand(c *gin.Context, chatID int64, messageId int) (string, string, error)
}

type telegramService struct {
	dbWordPress       *db.DB
	dbTracking        *db.DB
	wordPressPostRepo wprepositories.IFaWordpressPostRepository
	userTrackingRepo  tlRepositories.ITelegramChabotRepository
	chatResponseRepo  tlRepositories.ITelegramChabotResponseRepository
	redis             rediscmd.RedisCmd
}

type TelegramServiceArgs struct {
	dig.In
	DBWordPress       *db.DB `name:"faquizDB"`
	DBTracking        *db.DB `name:"trackingDB"`
	WordPressPostRepo wprepositories.IFaWordpressPostRepository
	UserTrackingRepo  tlRepositories.ITelegramChabotRepository
	ChatResponseRepo  tlRepositories.ITelegramChabotResponseRepository
	Redis             rediscmd.RedisCmd
}

func NewTelegramService(args TelegramServiceArgs) ITelegramService {
	return &telegramService{
		dbWordPress:       args.DBWordPress,
		dbTracking:        args.DBTracking,
		wordPressPostRepo: args.WordPressPostRepo,
		userTrackingRepo:  args.UserTrackingRepo,
		chatResponseRepo:  args.ChatResponseRepo,
		redis:             args.Redis,
	}
}

func (_this *telegramService) CallWebhook(c *gin.Context) (*meta.BasicResponse, error) {
	commandToSite := map[string][]string{
		"/r":   {viper.GetString(cfg.GoogleSeachEngineIDR), "radiologykey", "radiologykey.com"},
		"/m":   {viper.GetString(cfg.GoogleSeachEngineIDM), "musculoskeletalkey", "musculoskeletalkey.com"},
		"/p":   {viper.GetString(cfg.GoogleSeachEngineIDP), "plasticsurgerykey", "plasticsurgerykey.com"},
		"/d":   {viper.GetString(cfg.GoogleSeachEngineIDD), "pocketdentistry", "pocketdentistry.com"},
		"/t":   {viper.GetString(cfg.GoogleSeachEngineIDT), "thoracickey", "thoracickey.com"},
		"/v":   {viper.GetString(cfg.GoogleSeachEngineIDV), "veteriankey", "veteriankey.com"},
		"/n":   {viper.GetString(cfg.GoogleSeachEngineIDN), "neupsykey", "neupsykey.com"},
		"/nu":  {viper.GetString(cfg.GoogleSeachEngineIDNU), "nursekey", "nursekey.com"},
		"/o":   {viper.GetString(cfg.GoogleSeachEngineIDO), "obgynkey", "obgynkey.com"},
		"/on":  {viper.GetString(cfg.GoogleSeachEngineIDON), "oncohemakey", "oncohemakey.com"},
		"/e":   {viper.GetString(cfg.GoogleSeachEngineIDE), "entokey", "entokey.com"},
		"/c":   {viper.GetString(cfg.GoogleSeachEngineIDC), "clemedicine", "clemedicine.com"},
		"/b":   {viper.GetString(cfg.GoogleSeachEngineIDB), "basicmedicalkey", "basicmedicalkey.com"},
		"/a":   {viper.GetString(cfg.GoogleSeachEngineIDA), "aneskey", "aneskey.com"},
		"/ab":  {viper.GetString(cfg.GoogleSeachEngineIDAB), "abdominalkey", "abdominalkey.com"},
		"/app": {viper.GetString(cfg.GoogleSeachEngineID), "clinicalpub", "https://clinicalpub.com/"},
		"/web": {viper.GetString(cfg.GoogleSeachEngineID), "web"},
	}

	icons := map[string]string{
		"0": "📕",
		"1": "📙",
		"2": "📗",
		"3": "📘",
		"4": "📓",
	}

	var globalUpdate dtos.InteractionContext
	var update dtos.TelegramUpdate
	var response string

	action := "SEARCH"
	if err := c.ShouldBindJSON(&update); err != nil {
		ginLogger.Gin(c).Errorf("Error decoding update: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding update"})
		return nil, err
	}
	// Determine the type of update
	if update.Message != nil {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID
		userMessage := update.Message.Text
		messageId := update.Message.MessageID
		globalUpdate.UserID = userID
		globalUpdate.ChatID = chatID
		configKey := fmt.Sprintf("user_telegram_config:%d", userID)

		var userConfig dtos.UserTelegramConfig

		_, err := _this.redis.Get(configKey, &userConfig)
		if err != nil {
			ginLogger.Gin(c).Infof("No existing mode found or error retrieving mode for userID %d, setting default to /web. Error: %v", userID, err)
			userConfig.Mode = "/web"

			if err := _this.redis.Set(configKey, userConfig, 0); err != nil {
				ginLogger.Gin(c).Errorf("Error storing default UserTelegramConfig in Redis: %v", err)
			}
		} else {
			ginLogger.Gin(c).Infof("Successfully retrieved mode for userID %d: %s", userID, userConfig.Mode)
		}

		if !strings.HasPrefix(userMessage, "/") && userConfig.Mode == "/app" {
			userMessage = userConfig.Mode + " " + userMessage
		}
		globalUpdate.UserMessage = userMessage

		temp := strings.Fields(userMessage)[0]

		if value, exist := commandToSite[userMessage]; exist {
			userConfig := dtos.UserTelegramConfig{
				Mode: userMessage,
			}

			err := _this.redis.Set(configKey, userConfig, 0)
			if err != nil {
				ginLogger.Gin(c).Errorf("Error storing UserTelegramConfig in Redis: %v", err)
			}

			ginLogger.Gin(c).Infof("Changed search mode successfully")
			response, action, _ = _this.ModeCommand(c, chatID, messageId, value[1])

		} else if temp == "/start" {
			response, action, _ = _this.StartCommand(c, chatID, messageId)
		} else if temp == "/contact" {
			response, action, _ = _this.ContactCommand(c, chatID, messageId)
		} else if temp == "/specialty" {
			response, action, _ = _this.SpecialtiesCommand(c, chatID, messageId)
		} else {
			value := commandToSite[userConfig.Mode]

			replyText := ""
			var results []dtos.SearchResult

			results, _ = _this.PerformSearch(c, userMessage, value[0], temp == "/app", 1)
			if len(results) == 0 {
				emptyMessage := "Sorry! But we can not find any articles that matches your query! Please try another query!"
				responseTele, err := _this.SendMessage(c, chatID, emptyMessage, messageId)
				if err != nil {
					ginLogger.Gin(c).Errorf("Error sending response messagev: %v", err)
					return nil, errors.NewCusErr(errors.ErrCommonInvalidRequest)
				}
				response = responseTele
			} else {
				for index, item := range results {
					if userConfig.Mode == "/web" {
						cleanedURL := strings.Replace(item.DisplayLink, ".com", "", -1)
						replyText += fmt.Sprintf("<b>%s %s</b> \n 🌱<a href=\"%s\"><i><u>%s</u></i></a>    🔗 <a href=\"%s\"><b><u>View</u></b></a> \n\n", icons[strconv.Itoa(index)], item.Title, item.DisplayLink, cleanedURL, item.Link)
					} else {
						replyText += fmt.Sprintf("<b>%s %s</b> \n 🌱<a href=\"%s\"><i><u>%s</u></i></a>    🔗 <a href=\"%s\"><b><u>View</u></b></a> \n\n", icons[strconv.Itoa(index)], item.Title, value[2], value[1], item.Link)
					}
				}

				ginLogger.Gin(c).Infof("Message from user ID %d in chat ID %d: %s\n\n", userID, chatID, update.Message.Text)
				responseTele, err := _this.SendMessageWithCallBack(c, chatID, replyText, messageId, 1, globalUpdate)
				if err != nil {
					ginLogger.Gin(c).Errorf("Error sending response messagev: %v", err)
					return nil, errors.NewCusErr(errors.ErrCommonInvalidRequest)
				}
				response = responseTele
			}

		}

		var msg dtos.TelegramMessage
		if err := json.Unmarshal([]byte(response), &msg); err != nil {
			ginLogger.Gin(c).Errorf("Error parsing JSON: %v", err)
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
		_, err = _this.userTrackingRepo.Create(_this.dbTracking, userInfo)
		if err != nil {
			return nil, err
		}
	} else if update.CallbackQuery != nil {

		callbackData := update.CallbackQuery.Data
		var globalMessage dtos.InteractionContext
		callback, step, globalMessage := parseCallbackData(c, callbackData)

		ginLogger.Gin(c).Infof("Callback query from user ID %d at page %d: %s\n", globalMessage.UserID, step, callbackData)
		newCallbackData := fmt.Sprintf("%d_%s_%d", globalMessage.UserID, globalMessage.UserMessage, globalMessage.ChatID)
		switch callback {
		case "next":
			_this.HandleNextPage(c, globalMessage, update.CallbackQuery.Message.MessageID, step, newCallbackData)
		case "previous":
			_this.HandlePreviousPage(c, globalMessage, update.CallbackQuery.Message.MessageID, step, newCallbackData)
		}

		_this.AcknowledgeCallbackQuery(c, update.CallbackQuery.ID)
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

// AcknowledgeCallbackQuery sends an empty response to Telegram to stop the loading spinner on the button
func (_this *telegramService) AcknowledgeCallbackQuery(c *gin.Context, callbackQueryID string) {
	botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	request_url := fmt.Sprintf("https://api.telegram.org/bot%s/answerCallbackQuery", botToken)
	payload := map[string]string{"callback_query_id": callbackQueryID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error marshaling payload: %v", err)
		return
	}
	resp, err := http.Post(request_url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		ginLogger.Gin(c).Errorf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()
}

func (_this *telegramService) HandleNextPage(c *gin.Context, globalMessage dtos.InteractionContext, messageID int, currentPage int, callbackData string) {
	currentPage++

	if currentPage > 5 {
		ginLogger.Gin(c).Infof("Already at the last page, cannot go to next page.")
		return
	}
	_this.UpdateMessage(c, globalMessage, messageID, currentPage, callbackData)
}

func (_this *telegramService) HandlePreviousPage(c *gin.Context, globalMessage dtos.InteractionContext, messageID int, currentPage int, callbackData string) {
	currentPage--
	if currentPage < 1 {
		ginLogger.Gin(c).Infof("Already at the first page, cannot go to previous page.")
		return
	}

	_this.UpdateMessage(c, globalMessage, messageID, currentPage, callbackData)
}

// UpdateMessage edits a message's text
func (_this *telegramService) UpdateMessage(c *gin.Context, globalMessage dtos.InteractionContext, messageID int, currentPage int, callbackData string) {
	commandToSite := map[string][]string{
		"/r":   {viper.GetString(cfg.GoogleSeachEngineIDR), "radiologykey", "radiologykey.com"},
		"/m":   {viper.GetString(cfg.GoogleSeachEngineIDM), "musculoskeletalkey", "musculoskeletalkey.com"},
		"/p":   {viper.GetString(cfg.GoogleSeachEngineIDP), "plasticsurgerykey", "plasticsurgerykey.com"},
		"/d":   {viper.GetString(cfg.GoogleSeachEngineIDD), "pocketdentistry", "pocketdentistry.com"},
		"/t":   {viper.GetString(cfg.GoogleSeachEngineIDT), "thoracickey", "thoracickey.com"},
		"/v":   {viper.GetString(cfg.GoogleSeachEngineIDV), "veteriankey", "veteriankey.com"},
		"/n":   {viper.GetString(cfg.GoogleSeachEngineIDN), "neupsykey", "neupsykey.com"},
		"/nu":  {viper.GetString(cfg.GoogleSeachEngineIDNU), "nursekey", "nursekey.com"},
		"/o":   {viper.GetString(cfg.GoogleSeachEngineIDO), "obgynkey", "obgynkey.com"},
		"/on":  {viper.GetString(cfg.GoogleSeachEngineIDON), "oncohemakey", "oncohemakey.com"},
		"/e":   {viper.GetString(cfg.GoogleSeachEngineIDE), "entokey", "entokey.com"},
		"/c":   {viper.GetString(cfg.GoogleSeachEngineIDC), "clemedicine", "clemedicine.com"},
		"/b":   {viper.GetString(cfg.GoogleSeachEngineIDB), "basicmedicalkey", "basicmedicalkey.com"},
		"/a":   {viper.GetString(cfg.GoogleSeachEngineIDA), "aneskey", "aneskey.com"},
		"/ab":  {viper.GetString(cfg.GoogleSeachEngineIDAB), "abdominalkey", "abdominalkey.com"},
		"/app": {viper.GetString(cfg.GoogleSeachEngineID), "clinicalpub", "https://clinicalpub.com/"},
		"/web": {viper.GetString(cfg.GoogleSeachEngineID), "web", ""},
	}

	icons := map[string]string{
		"0": "📕",
		"1": "📙",
		"2": "📗",
		"3": "📘",
		"4": "📓",
	}

	botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	editUrl := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", botToken)

	temp := strings.Fields(globalMessage.UserMessage)[0]
	replyText := ""
	if value, exists := commandToSite[temp]; exists {
		_, message := TakeQueryAndAction(globalMessage.UserMessage)

		results, _ := _this.PerformSearch(c, message, value[0], temp == "/app", currentPage*5-4)

		for index, item := range results {
			replyText += fmt.Sprintf("<b>%s %s</b> \n 🌱<a href=\"%s\"><i><u>%s</u></i></a>    🔗 <a href=\"%s\"><b><u>View</u></b></a> \n\n", icons[strconv.Itoa(index)], item.Title, value[2], value[1], item.Link)
		}
	} else {
		value = commandToSite["/web"]

		results, _ := _this.PerformSearch(c, globalMessage.UserMessage, value[0], false, currentPage*5-4)

		for index, item := range results {
			cleanedURL := strings.Replace(item.DisplayLink, ".com", "", -1)
			replyText += fmt.Sprintf("<b>%s %s</b> \n 🌱<a href=\"%s\"><i><u>%s</u></i></a>    🔗 <a href=\"%s\"><b><u>View</u></b></a> \n\n", icons[strconv.Itoa(index)], item.Title, item.DisplayLink, cleanedURL, item.Link)
		}
	}
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "Previous", "callback_data": fmt.Sprintf("previous_%d_%s", currentPage, callbackData)},
				{"text": "Next", "callback_data": fmt.Sprintf("next_%d_%s", currentPage, callbackData)},
			},
		},
	}

	payload := map[string]interface{}{
		"chat_id":      globalMessage.ChatID,
		"message_id":   messageID,
		"text":         replyText,
		"parse_mode":   "HTML",
		"reply_markup": keyboard,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		// Handle error
		ginLogger.Gin(c).Errorf("Error marshaling payload: %v", err)
		return
	}
	resp, err := http.Post(editUrl, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		// Handle error
		ginLogger.Gin(c).Errorf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()
	// Check response status code, etc.
}
func (_this *telegramService) SendMessageWithCallBack(c *gin.Context, chatID int64, message string, messageId int, currentPage int, globalUpdate dtos.InteractionContext) (string, error) {
	// API endpoint for sending messages
	botToken := viper.GetString(cfg.ConfigTelegramBotToken)
	sendMessageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	globalMessageString := fmt.Sprintf("%d_%s_%d", globalUpdate.UserID, globalUpdate.UserMessage, globalUpdate.ChatID)

	inlineKeyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{"text": "Next", "callback_data": fmt.Sprintf("next_%d_%s", currentPage, globalMessageString)},
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
	ginLogger.Gin(c).Infof("Response from Telegram: %s\n", string(responseBody))

	return string(responseBody), nil
}

// SendMessage sends a message to a user in Telegram
func (_this *telegramService) SendMessage(c *gin.Context, chatID int64, message string, messageId int) (string, error) {
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

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response from Telegram:", err)
		return "", err
	}
	ginLogger.Gin(c).Errorf("Response from Telegram: %s\n", string(responseBody))

	return string(responseBody), nil
}

//func (_this *telegramService) FeedbackCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
//	action := "REPORT"
//	feedbackMessage := "We have received you feedback!"
//	responseTele, err := _this.SendMessage(c, chatID, feedbackMessage, messageId)
//	if err != nil {
//		ginLogger.Gin(c).Errorf("Error sending response messagev: %", err)
//		return "", "ERROR", errors.NewCusErr(errors.ErrCommonInvalidRequest)
//	}
//	return responseTele, action, nil
//}

func (_this *telegramService) SpecialtiesCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
	specialtyMessage, _ := _this.chatResponseRepo.FindByName(_this.dbTracking, "SPECIALTY")
	responseTele, err := _this.SendMessage(c, chatID, specialtyMessage, messageId)
	action := "ABOUT"
	if err != nil {
		ginLogger.Gin(c).Errorf("Error sending response messagev:  %v", err)
		return "", "ERROR", errors.NewCusErr(errors.ErrCommonInvalidRequest)
	}
	return responseTele, action, nil
}

func (_this *telegramService) ContactCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
	aboutMessage, _ := _this.chatResponseRepo.FindByName(_this.dbTracking, "CONTACT")
	responseTele, err := _this.SendMessage(c, chatID, aboutMessage, messageId)
	action := "CONTACT"
	if err != nil {
		ginLogger.Gin(c).Errorf("Error sending response messagev:  %v", err)
		return "", "ERROR", errors.NewCusErr(errors.ErrCommonInvalidRequest)
	}
	return responseTele, action, nil
}

func (_this *telegramService) StartCommand(c *gin.Context, chatID int64, messageId int) (string, string, error) {
	welcomeMessage, _ := _this.chatResponseRepo.FindByName(_this.dbTracking, "START")
	responseTele, err := _this.SendMessage(c, chatID, welcomeMessage, messageId)
	action := "START"
	if err != nil {
		ginLogger.Gin(c).Errorf("Error sending response messagev: %v", err)
		return "", "ERROR", errors.NewCusErr(errors.ErrCommonInvalidRequest)
	}
	return responseTele, action, nil
}

func (_this *telegramService) ModeCommand(c *gin.Context, chatID int64, messageId int, mode string) (string, string, error) {
	modeMessage, _ := _this.chatResponseRepo.FindByName(_this.dbTracking, "MODE")
	responseTele, err := _this.SendMessage(c, chatID, modeMessage+mode, messageId)
	action := "CONTACT"
	if err != nil {
		ginLogger.Gin(c).Errorf("Error sending response messagev:  %v", err)
		return "", "ERROR", errors.NewCusErr(errors.ErrCommonInvalidRequest)
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

func (_this *telegramService) PerformSearch(c *gin.Context, query string, site string, isApp bool, start int) ([]dtos.SearchResult, error) {
	res := make([]dtos.SearchResult, 0)
	apiKey := getNextApiKey()
	if isApp == true {
		posts, err := _this.wordPressPostRepo.GetPostsByTitle(_this.dbWordPress, query, start)
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
		searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%d&start=%d",
			apiKey, site, url.QueryEscape(query), left, start)

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

func parseCallbackData(c *gin.Context, callbackData string) (string, int, dtos.InteractionContext) {
	var globalValue dtos.InteractionContext

	parts := strings.Split(callbackData, "_")
	action := parts[0]
	pageNumber, err := strconv.Atoi(parts[1])

	chatId, err := strconv.ParseInt(parts[4], 10, 64)
	if err == nil {
		ginLogger.Gin(c).Errorf("Err: %v", err)
	}

	userId, err := strconv.ParseInt(parts[2], 10, 64)
	if err == nil {
		ginLogger.Gin(c).Errorf("Err: %v", err)
	}
	globalValue.UserID = userId
	globalValue.ChatID = chatId
	globalValue.UserMessage = parts[3]
	return action, pageNumber, globalValue
}
