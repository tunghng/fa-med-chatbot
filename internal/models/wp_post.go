package models

import "time"

const TableNameWPPost = "wplw_posts"

type WPPost struct {
	ID                  int64     `gorm:"column:ID;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	PostAuthor          int64     `gorm:"column:post_author" json:"postAuthor"`
	PostDate            time.Time `gorm:"column:post_date" json:"postDate"`
	PostDateGmt         time.Time `gorm:"column:post_date_gmt" json:"postDateGmt"`
	PostContent         string    `gorm:"column:post_content" json:"postContent"`
	PostTitle           string    `gorm:"column:post_title" json:"postTitle"`
	PostExcerpt         string    `gorm:"column:post_excerpt" json:"postExcerpt"`
	PostStatus          string    `gorm:"column:post_status" json:"postStatus"`
	CommentStatus       string    `gorm:"column:comment_status" json:"commentStatus"`
	PingStatus          string    `gorm:"column:ping_status" json:"pingStatus"`
	PostPassword        string    `gorm:"column:post_password" json:"postPassword"`
	PostName            string    `gorm:"column:post_name" json:"postName"`
	ToPing              string    `gorm:"column:to_ping" json:"toPing"`
	Pinged              string    `gorm:"column:pinged" json:"pinged"`
	PostModified        time.Time `gorm:"column:post_modified" json:"postModified"`
	PostModifiedGmt     time.Time `gorm:"column:post_modified_gmt" json:"postModifiedGmt"`
	PostContentFiltered string    `gorm:"column:post_content_filtered" json:"postContentFiltered"`
	PostParent          int64     `gorm:"column:post_parent" json:"postParent"`
	GUID                string    `gorm:"column:guid" json:"guid"`
	MenuOrder           int64     `gorm:"column:menu_order" json:"menuOrder"`
	PostType            string    `gorm:"column:post_type" json:"postType"`
	PostMimeType        string    `gorm:"column:post_mime_type" json:"postMimeType"`
	CommentCount        int64     `gorm:"column:comment_count" json:"commentCount"`
}

func (WPPost) TableName() string {
	return TableNameWPPost // Replace with your actual table name
}
