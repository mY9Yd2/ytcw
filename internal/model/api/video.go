package model

import "time"

type VideoResponse struct {
	ID        uint           `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	FullTitle string         `json:"full_title"`
	DisplayID string         `json:"display_id"`
	Duration  uint           `json:"duration"`
	Language  *string        `json:"language"`
	Thumbnail string         `json:"thumbnail"`
	Channel   ChannelSummary `json:"channel"`
}
