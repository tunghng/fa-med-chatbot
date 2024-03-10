package telegram

import (
	"med-chat-bot/internal/models"
	"med-chat-bot/pkg/db"
)

type ITelegramChabotRepository interface {
	FindById(_db *db.DB, id int64) (*models.UserTrackingChatBot, error)
	Create(db *db.DB, obj *models.UserTrackingChatBot) (*models.UserTrackingChatBot, error)
	Update(db *db.DB, obj *models.UserTrackingChatBot) error
}

type dTeleRepo struct{}

func NewTelegramChabotRepository() ITelegramChabotRepository {
	return &dTeleRepo{}
}

func (_this *dTeleRepo) Create(db *db.DB, obj *models.UserTrackingChatBot) (*models.UserTrackingChatBot, error) {
	if err := db.DB().Table(models.TableNameTrackingUser).Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

func (_this *dTeleRepo) Update(db *db.DB, item *models.UserTrackingChatBot) error {
	return db.DB().Table(models.TableNameTrackingUser).
		Where("id = ? ", item.ID).
		Update(item).Error
}

func (_this *dTeleRepo) FindById(_db *db.DB, id int64) (*models.UserTrackingChatBot, error) {
	var item models.UserTrackingChatBot

	dbTnx := _db.DB().Table(models.TableNameTrackingUser)

	err := dbTnx.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}
