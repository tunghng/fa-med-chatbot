package wordpress

import (
	"med-chat-bot/internal/models"
	"med-chat-bot/pkg/db"
	"strings"
)

type IFaWordpressPostRepository interface {
	Create(db *db.DB, obj *models.WPPost) (*models.WPPost, error)
	Update(db *db.DB, item *models.WPPost) error
	GetPostsByTitle(db *db.DB, title string) ([]models.WPPost, error)
}

type wordpressPostRepo struct{}

func NewWordpressPostRepository() IFaWordpressPostRepository {
	return &wordpressPostRepo{}
}

func (_this *wordpressPostRepo) Create(db *db.DB, obj *models.WPPost) (*models.WPPost, error) {
	if err := db.DB().Table(models.TableNameWPPost).Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

func (_this *wordpressPostRepo) Update(db *db.DB, item *models.WPPost) error {
	return db.DB().Table(models.TableNameWPPost).
		Where("ID = ? ", item.ID).
		Update(item).Error
}

func (_this *wordpressPostRepo) GetPostsByTitle(db *db.DB, title string) ([]models.WPPost, error) {
	var posts []models.WPPost
	processed := strings.TrimSpace(title[7:])
	err := db.DB().Table("wplw_posts").
		Select("post_title AS title, guid AS link").
		Where("post_title LIKE ?", "%"+processed+"%").
		Where("post_status = ?", "publish").
		Limit(2).
		Find(&posts).Error

	return posts, err
}
