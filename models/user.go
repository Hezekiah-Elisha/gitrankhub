package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex;not null" json:"username"`
	Name      string `gorm:"not null" json:"name"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `json:"-"`
	Role      string `gorm:"not null;default:'user'" json:"role"`
	AvatarURL string `gorm:"not null" json:"avatar_url"`
	Bio       string `gorm:"not null" json:"bio"`
}
