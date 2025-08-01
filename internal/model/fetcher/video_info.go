package model

import "github.com/mY9Yd2/ytcw/internal/model"

type VideoInfo struct {
	ChannelInfo
	Timestamp int64  `json:"timestamp"`
	FullTitle string `json:"fulltitle"`
	DisplayID string `json:"display_id"`
	Duration  uint   `json:"duration"`
	Language  string `json:"language"`
	Thumbnail string
	VideoType model.VideoType
}
