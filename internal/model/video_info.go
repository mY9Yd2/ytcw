package model

type VideoInfo struct {
	ChannelInfo
	Timestamp int64  `json:"timestamp"`
	FullTitle string `json:"fulltitle"`
	DisplayID string `json:"display_id"`
	Duration  int    `json:"duration"`
	Language  string `json:"language"`
	Thumbnail string
}
