package elasticsearch

type ItemsDto struct {
	Id             string  `json:"id"`
	CreatedAt      int64   `json:"createdAt"`
	ItemId         int64   `json:"itemId"`
	ItemType       string  `json:"itemType"`
	CoverUrl       *string `json:"coverUrl"`
	Description    string  `json:"description"`
	SubDescription string  `json:"subDescription"`
	Status         string  `json:"status"`
	AuthorName     string  `json:"authorName"`
	AuthorId       int64   `json:"authorId"`
	AuthorAvatar   string  `json:"authorAvatar"`
	TotalCards     int64   `json:"totalCards"`
}

type ItemsAutoDto struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

const (
	ElkTypeAuthor  = "AUTHOR"
	ElkTypeTitle   = "TITLE"
	ElkTypeContent = "CONTENT"
	ElkTypeImage   = "IMAGE"
)
