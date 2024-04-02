package wordpress

import (
	"med-chat-bot/internal/models"
	"med-chat-bot/pkg/db"
	"strings"
)

type IFaWordpressPostRepository interface {
	Create(db *db.DB, obj *models.WPPost) (*models.WPPost, error)
	Update(db *db.DB, item *models.WPPost) error
	GetPostsByTitle(db *db.DB, title string, start int) ([]models.WPPost, error)
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

func (_this *wordpressPostRepo) GetPostsByTitle(db *db.DB, title string, start int) ([]models.WPPost, error) {
	var posts []models.WPPost
	processed := strings.TrimSpace(title)

	err := db.DB().Table(models.TableNameWPPost).
		Select("post_title, guid"). // Fetch columns directly as named in the struct
		Where("LOWER(post_title) LIKE ?", "%"+strings.ToLower(processed)+"%").
		Where("post_status = ?", "publish").
		Limit(5).
		Offset(start).
		Find(&posts).Error

	return posts, err
}
