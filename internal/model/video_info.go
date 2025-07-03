package model

type VideoInfo struct {
	UploaderID string `json:"uploader_id"`
	ChannelID  string `json:"channel_id"`
	Channel    string `json:"channel"`
	Timestamp  int64  `json:"timestamp"`
	FullTitle  string `json:"fulltitle"`
	DisplayID  string `json:"display_id"`
	Duration   int    `json:"duration"`
	Language   string `json:"language"`
	Thumbnail  string
}
