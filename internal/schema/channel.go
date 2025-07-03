package schema

import (
	"gorm.io/gorm"
	"time"
)

type Channel struct {
	gorm.Model
	Videos     []Video `gorm:"foreignKey:ChannelRefer;references:ID;"`
	CategoryID *uint
	Category   *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	UploaderID string     `gorm:"unique;size:50;"`
	ChannelID  string     `gorm:"unique;size:30;"`
	Channel    string     `gorm:"size:80;"`
	LastFetch  *time.Time `gorm:"index"`
}
