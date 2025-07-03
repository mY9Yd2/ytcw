package mapper

import (
	"time"
	"ytcw/internal/model"
	"ytcw/internal/schema"
)

func MapVideoInfoToVideo(info model.VideoInfo) schema.Video {
	return schema.Video{
		Timestamp: time.Unix(info.Timestamp, 0),
		FullTitle: info.FullTitle,
		DisplayID: info.DisplayID,
		Duration:  info.Duration,
		Language:  &info.Language,
		Thumbnail: info.Thumbnail,
		Channel: schema.Channel{
			UploaderID: info.UploaderID,
			ChannelID:  info.ChannelID,
			Channel:    info.Channel,
		},
	}
}
