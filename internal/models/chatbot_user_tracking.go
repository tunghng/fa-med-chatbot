package models

import "time"

const TableNameTrackingUser = "user_tracking_information"

type UserTrackingChatBot struct {
	ID          int64     `gorm:"column:ID;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	UserID      int64     `gorm:"column:user_id" json:"user_id"`
	ChatID      int64     `gorm:"column:chat_id" json:"chat_id"`
	Action      string    `gorm:"column:action" json:"action"`
	Username    string    `gorm:"column:username" json:"username"`
	MessageText string    `gorm:"column:message_text" json:"message_text"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	MetaData    string    `gorm:"column:meta_data" json:"meta_data"`
}

func (UserTrackingChatBot) TableName() string {
	return TableNameTrackingUser // Replace with your actual table name
}
