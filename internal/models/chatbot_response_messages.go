package models

import "time"

const TableNameChatbotResponse = "chatbot_response_message"

type ChatbotResponse struct {
	ID          int64     `gorm:"column:ID;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	MessageText string    `gorm:"column:message_text" json:"messageText"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
	ActionName  string    `gorm:"column:action_name" json:"actionName"`
}

func (ChatbotResponse) TableName() string {
	return TableNameChatbotResponse
}
