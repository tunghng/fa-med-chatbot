package dtos

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
	}
}
