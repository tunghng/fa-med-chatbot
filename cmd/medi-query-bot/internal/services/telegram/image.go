package telegram

import (
	"bytes"
	"encoding/json"
	"med-chat-bot/internal/dtos"
	"med-chat-bot/internal/errors"
	"med-chat-bot/internal/ginLogger"
	"med-chat-bot/internal/meta"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type IImageService interface {
	UploadImageAsDocumentToTelegram(c *gin.Context) (*meta.BasicResponse, error)
	GetImageFromTelegram(c *gin.Context, fileID string) (*meta.BasicResponse, error)
}

type imageService struct {
}

type ImageServiceArgs struct {
	dig.In
}

func NewImageService(args ImageServiceArgs) IImageService {
	return &imageService{}
}

// Test variables
const (
	filePath = "/Users/htungg/Downloads/Cat03.jpg"
	botToken = "6723992316:AAHJf3HxVBXyocUIPQ2t2KJHM1SmsU1ZrSg"
	chatID   = "1148785584"
	savePath = "/Users/htungg/Downloads/image1.jpg"
)

func (_this imageService) UploadImageAsDocumentToTelegram(c *gin.Context) (*meta.BasicResponse, error) {
	// Open file from local path ("/Users/htungg/Downloads/Cat03.jpg")
	file, err := os.Open(filePath)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	defer file.Close()

	// Turn file into document
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("document", filePath)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}

	_ = writer.WriteField("chat_id", chatID)
	writer.Close()

	// HTTP Request to Telegram
	url := "https://api.telegram.org/bot" + botToken + "/sendDocument"
	resp, err := http.Post(url, writer.FormDataContentType(), body)

	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}

	// Parse Response to dto
	var r dtos.TelegramResponse
	if err := json.Unmarshal(responseBody, &r); err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}

	if !r.OK {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}

	response := &meta.BasicResponse{
		Meta: meta.Meta{
			Code:    http.StatusOK,
			Message: "",
		},
		Data: r,
	}

	//fileID := r.Result.Document.FileID
	return response, nil
}

func (_this imageService) GetImageFromTelegram(c *gin.Context, fileID string) (*meta.BasicResponse, error) {
	getFileURL := "https://api.telegram.org/bot" + botToken + "/getFile?file_id=" + fileID

	resp, err := http.Get(getFileURL)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var getFileResponse dtos.GetFileResponse

	if err := json.NewDecoder(resp.Body).Decode(&getFileResponse); err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, err
	}

	if !getFileResponse.OK {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, err
	}

	downloadURL := "https://api.telegram.org/file/bot" + botToken + "/" + getFileResponse.Result.FilePath
	fileResp, err := http.Get(downloadURL)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	defer fileResp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		ginLogger.Gin(c).Errorf("Error: %v", err)
		return nil, errors.NewCusErr(errors.ErrCommonInternalServer)
	}
	defer out.Close()

	_, err = io.Copy(out, fileResp.Body)

	response := &meta.BasicResponse{
		Meta: meta.Meta{
			Code:    http.StatusOK,
			Message: "",
		},
		Data: nil,
	}

	return response, err
}
