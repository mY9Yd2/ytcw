package schema

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	UUIDModel
	Videos        []Video    `gorm:"foreignKey:ChannelRefer;references:ID;"`
	CategoryRefer *uuid.UUID `gorm:"type:uuid;"`
	Category      *Category  `gorm:"foreignKey:CategoryRefer;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	UploaderID string     `gorm:"unique;size:50;"`
	ChannelID  string     `gorm:"unique;size:30;"`
	Channel    string     `gorm:"size:80;"`
	LastFetch  *time.Time `gorm:"index"`
	DisabledAt *time.Time `gorm:"index"`
}
