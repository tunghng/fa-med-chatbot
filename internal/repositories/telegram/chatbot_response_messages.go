package telegram

import (
	"med-chat-bot/internal/models"
	"med-chat-bot/pkg/db"
)

type ITelegramChabotResponseRepository interface {
	FindByName(_db *db.DB, actionName string) (string, error)
}

type dTeleResRepo struct{}

func NewTelegramChabotResponseRepository() ITelegramChabotResponseRepository {
	return &dTeleResRepo{}
}

func (_this *dTeleResRepo) FindByName(_db *db.DB, actionName string) (string, error) {
	var item models.ChatbotResponse

	dbTnx := _db.DB().Table(models.TableNameChatbotResponse)

	err := dbTnx.Where("action_name = ?", actionName).First(&item).Error
	if err != nil {
		return "", err
	}

	return item.MessageText, nil
}
