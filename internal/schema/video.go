package schema

import (
	"github.com/google/uuid"
	"github.com/mY9Yd2/ytcw/internal/model"
	"time"
)

type Video struct {
	UUIDModel
	ChannelRefer uuid.UUID `gorm:"type:uuid;"`
	Channel      Channel   `gorm:"foreignKey:ChannelRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Timestamp time.Time
	FullTitle string
	DisplayID string `gorm:"unique;size:20;"`
	Duration  uint
	Language  *string         `gorm:"size:6;"`
	Thumbnail string          `gorm:"size:14;"`
	VideoType model.VideoType `gorm:"size:7;"`
}
