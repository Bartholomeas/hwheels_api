package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New().String()
	return
}
