package model

import (
	"github.com/google/uuid"
	"github.com/mY9Yd2/ytcw/internal/model"
	"time"
)

type VideoResponse struct {
	ID        uuid.UUID       `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	FullTitle string          `json:"full_title"`
	DisplayID string          `json:"display_id"`
	Duration  uint            `json:"duration"`
	Language  *string         `json:"language"`
	Thumbnail string          `json:"thumbnail"`
	Channel   ChannelSummary  `json:"channel"`
	VideoType model.VideoType `json:"video_type"`
}
