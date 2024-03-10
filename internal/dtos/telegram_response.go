package dtos

// TelegramUpdate represents the structure of an update received from Telegram.
type TelegramUpdate struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
	// Add other fields as needed.
}

// Message represents a Telegram message.
type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text"`
}

type User struct {
	ID        int64  `json:"id"` // User ID
	FirstName string `json:"first_name"`
	// Include other fields as needed
}

type Chat struct {
	ID int64 `json:"id"` // Chat ID
	// Include other fields as needed
}

// CallbackQuery represents a callback query from an inline keyboard.
type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message,omitempty"`
	Data    string   `json:"data"`
}

type SearchResult struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type GoogleSearchResponse struct {
	Items []SearchResult `json:"items"`
}

type LinkConversionRequest struct {
	WorkspaceID string `json:"workspace_id"`
	Url         string `json:"url"`
}

type LinkConversionResponse struct {
	Link struct {
		FullURL string `json:"full_url"`
	} `json:"link"`
	Status string `json:"status"`
}

type MessageBody struct {
	ChatID int64  `json:"id"`
	Text   string `json:"message"`
}

type TelegramResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Document struct {
			FileID   string `json:"file_id"`
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
		} `json:"document"`
	} `json:"result"`
}

type GetFileResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		FilePath string `json:"file_path"`
	} `json:"result"`
}
