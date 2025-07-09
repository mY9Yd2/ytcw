package schema

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UUIDModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID, err = uuid.NewV7()
	}
	return
}
