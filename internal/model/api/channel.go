package model

import "time"

type ChannelResponse struct {
	ID         uint              `json:"id"`
	UploaderID string            `json:"uploader_id"`
	ChannelID  string            `json:"channel_id"`
	Channel    string            `json:"channel"`
	LastFetch  *time.Time        `json:"last_fetch"`
	DisabledAt *time.Time        `json:"disabled_at"`
	Category   *CategoryResponse `json:"category"`
}
