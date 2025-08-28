package content

import (
	"time"

	"github.com/google/uuid"
)

type ChannelResponse struct {
	ID            uuid.UUID         `json:"id"`
	UploaderID    string            `json:"uploader_id"`
	ChannelID     string            `json:"channel_id"`
	Channel       string            `json:"channel"`
	LastFetch     *time.Time        `json:"last_fetch"`
	DisabledAt    *time.Time        `json:"disabled_at"`
	DisabledUntil *time.Time        `json:"disabled_until"`
	Category      *CategoryResponse `json:"category"`
}

type ChannelSummary struct {
	ID       uuid.UUID         `json:"id"`
	Channel  string            `json:"channel"`
	Category *CategoryResponse `json:"category"`
}

type VideoResponse struct {
	ID        uuid.UUID      `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	FullTitle string         `json:"full_title"`
	DisplayID string         `json:"display_id"`
	Duration  uint           `json:"duration"`
	Language  *string        `json:"language"`
	Thumbnail string         `json:"thumbnail"`
	Channel   ChannelSummary `json:"channel"`
	VideoType VideoType      `json:"video_type"`
}

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
