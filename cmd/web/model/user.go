package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         uint
	Name       string `validate:"required"`
	Username   string `gorm:"uniqueIndex;not nul"`
	Email      string `gorm:"uniqueIndex;not null" validate:"required,email"`
	Password   string `gorm:"not null"`
	Bio        string
	Image      string
	Followers  []Follow  `gorm:"foreignKey:FollowerID"`
	Followings []Follow  `gorm:"foreignKey:FollowingID"`
	Favorites  []Article `gorm:"many2many:article_favorite;"`
}

type Follow struct {
	Follower    User
	FollowerID  uint `gorm:"column:user_id;primaryKey" sql:"type:int not null"`
	Following   User
	FollowingID uint `gorm:"column:follower_id;primaryKey" sql:"type:int not null"`
}

func (u *User) HashPassword() {
	h, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(h)
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

func (u User) FollowedBy(id uint) bool {
	if u.Followers == nil {
		return false
	}

	for _, f := range u.Followers {
		if f.FollowingID == id {
			return true
		}
	}

	return false
}

func (u User) FollowersCount() int {
	return len(u.Followers)
}

func (Follow) TableName() string {
	return "user_follower"
}
