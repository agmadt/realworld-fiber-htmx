package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Slug        string `gorm:"uniqueIndex;not null"`
	Title       string `gorm:"not null" validate:"required"`
	Description string `validate:"required"`
	Body        string `validate:"required"`
	User        User   `validate:"-"`
	UserID      uint
	Comments    []Comment
	Favorites   []User `gorm:"many2many:article_favorite;"`
	Tags        []Tag  `gorm:"many2many:article_tag;"`
	IsFavorited bool   `gorm:"-"`
}

func (Article Article) GetFormattedCreatedAt() string {
	dateLayout := "Jan 02, 2006"
	return Article.CreatedAt.Format(dateLayout)
}

func (Article Article) GetFavoriteCount() int {
	return len(Article.Favorites)
}

func (Article Article) FavoritedBy(id uint) bool {
	if Article.Favorites == nil {
		return false
	}

	for _, u := range Article.Favorites {
		if u.ID == id {
			return true
		}
	}

	return false
}

func (Article Article) GetTagsAsCommaSeparated() string {
	tagsText := ""

	for i := 0; i < len(Article.Tags); i++ {
		tagsText += Article.Tags[i].Name + ","
	}

	return tagsText
}
