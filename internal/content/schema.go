package content

import (
	"time"

	"github.com/google/uuid"
	"github.com/mY9Yd2/ytcw/internal/common"
)

type Channel struct {
	common.UUIDModel
	Videos        []Video    `gorm:"foreignKey:ChannelRefer;references:ID;"`
	CategoryRefer *uuid.UUID `gorm:"type:uuid;"`
	Category      *Category  `gorm:"foreignKey:CategoryRefer;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	UploaderID    string     `gorm:"unique;size:50;"`
	ChannelID     string     `gorm:"unique;size:30;"`
	Channel       string     `gorm:"size:80;"`
	LastFetch     *time.Time `gorm:"index"`
	DisabledAt    *time.Time `gorm:"index"`
	DisabledUntil *time.Time `gorm:"index"`
}

type Video struct {
	common.UUIDModel
	ChannelRefer uuid.UUID `gorm:"type:uuid;"`
	Channel      Channel   `gorm:"foreignKey:ChannelRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Timestamp time.Time
	FullTitle string
	DisplayID string `gorm:"unique;size:20;"`
	Duration  uint
	Language  *string   `gorm:"size:6;"`
	Thumbnail string    `gorm:"size:14;"`
	VideoType VideoType `gorm:"size:7;"`
}

type Category struct {
	common.UUIDModel
	Channels []Channel `gorm:"foreignKey:CategoryRefer;"`

	Name string `gorm:"unique;size:40;"`
}
