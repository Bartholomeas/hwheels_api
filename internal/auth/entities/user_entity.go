package entities

import (
	"github.com/bartholomeas/hwheels_api/internal/common/entities"
	userEntities "github.com/bartholomeas/hwheels_api/internal/user/entities"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	entities.Base
	Email    string                    `json:"email" gorm:"uniqueIndex;not null"`
	Username string                    `json:"username" gorm:"uniqueIndex;not null;size:255"`
	Password string                    `json:"password" gorm:"not null"`
	Role     Role                      `json:"role" gorm:"type:varchar(20);default:'user'"`
	Profile  *userEntities.UserProfile `json:"profile" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "users"
}
