package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Number    string         `json:"number" gorm:"type:varchar(14);uniqueIndex;not null"`
	Type      string         `json:"type" gorm:"type:varchar(4);not null"`
	Blocked   bool           `json:"blocked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

