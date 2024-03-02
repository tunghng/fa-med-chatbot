package repositories

import (
	"med-chat-bot/pkg/db"
)

type Post struct {
	Title string
	Link  string
}

type ILinkRepository interface {
	GetPostsByTitle(db *db.DB, name string) ([]Post, error)
}

type linkRepo struct{}

func NewLinkRepository() ILinkRepository { return &linkRepo{} }

func (_this *linkRepo) GetPostsByTitle(db *db.DB, name string) ([]Post, error) {
	var posts []Post
	err := db.DB().Table("wplw_posts").
		Select("post_title AS title", "guid AS link").
		Where("post_title LIKE ?", "%"+name+"%").
		Where("post_status = ?", "publish").
		Limit(5).
		Find(&posts).Error

	return posts, err
}
