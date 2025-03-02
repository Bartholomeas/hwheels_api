package entities

import "github.com/bartholomeas/hwheels_api/api/common/entities"

type UserProfile struct {
	entities.Base
	FirstName *string `json:"first_name,omitempty" gorm:"size:255"`
	LastName  *string `json:"last_name,omitempty" gorm:"size:255"`
	Phone     *string `json:"phone,omitempty" gorm:"size:255"`
	AvatarURL *string `json:"avatar_url,omitempty" gorm:"size:255"`
	UserID    string  `json:"user_id" gorm:"type:uuid;not null;unique"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
