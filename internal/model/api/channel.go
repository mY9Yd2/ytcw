package model

import (
	"github.com/google/uuid"
	"time"
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
