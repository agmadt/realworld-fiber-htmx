package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Article   Article `validate:"-"`
	ArticleID uint
	User      User `validate:"-"`
	UserID    uint
	Body      string `validate:"required"`
}

func (Comment Comment) GetFormattedCreatedAt() string {
	dateLayout := "Jan 02, 2006"
	return Comment.CreatedAt.Format(dateLayout)
}
