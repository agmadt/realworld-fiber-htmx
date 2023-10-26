package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name     string    `gorm:"uniqueIndex"`
	Articles []Article `gorm:"many2many:article_tag;"`
}
