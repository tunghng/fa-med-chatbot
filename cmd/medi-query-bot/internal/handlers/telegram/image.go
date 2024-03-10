package telegram

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"med-chat-bot/cmd/medi-query-bot/internal/services/telegram"
	"med-chat-bot/internal/handlers"
	"strings"
)

type ImageHandler struct {
	handlers.BaseHandler
	imageService telegram.IImageService
}

type ImageHandlerParams struct {
	dig.In
	BaseHandler  handlers.BaseHandler
	ImageService telegram.IImageService
}

func NewImageHandler(params ImageHandlerParams) *ImageHandler {
	return &ImageHandler{
		BaseHandler:  params.BaseHandler,
		imageService: params.ImageService,
	}
}

func (_this *ImageHandler) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := _this.imageService.UploadImageAsDocumentToTelegram(c)
		_this.HandleResponse(c, resp, err)
	}
}

// DownloadImage use test fileId gotten from upload image as query param
func (_this *ImageHandler) DownloadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileId := strings.TrimSpace(c.Query("fileId"))
		resp, err := _this.imageService.GetImageFromTelegram(c, fileId)
		_this.HandleResponse(c, resp, err)
	}
}
