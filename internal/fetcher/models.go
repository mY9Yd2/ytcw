package fetcher

import "github.com/mY9Yd2/ytcw/internal/content"

type ChannelInfo struct {
	UploaderID string `json:"uploader_id"`
	ChannelID  string `json:"channel_id"`
	Channel    string `json:"channel"`
}

type VideoInfo struct {
	ChannelInfo
	Timestamp int64  `json:"timestamp"`
	FullTitle string `json:"fulltitle"`
	DisplayID string `json:"display_id"`
	Duration  uint   `json:"duration"`
	Language  string `json:"language"`
	Thumbnail string
	VideoType content.VideoType
}
