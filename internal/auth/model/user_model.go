package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string `json:"id" gorm:"primary_key"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"` // TODO: There should be enum somehow
	CreatedAt time.Time
	UpdatedAt time.Time
}
