package schema

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	ChannelRefer uint
	Channel      Channel `gorm:"foreignKey:ChannelRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Timestamp time.Time
	FullTitle string
	DisplayID string `gorm:"unique;size:20;"`
	Duration  int
	Language  *string `gorm:"size:6;"`
	Thumbnail string  `gorm:"size:14;"`
}
