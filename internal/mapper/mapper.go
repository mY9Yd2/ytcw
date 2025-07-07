package mapper

import (
	model "github.com/mY9Yd2/ytcw/internal/model/fetcher"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"time"
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

func MapChannelInfoToChannel(info model.ChannelInfo) schema.Channel {
	return schema.Channel{
		UploaderID: info.UploaderID,
		ChannelID:  info.ChannelID,
		Channel:    info.Channel,
	}
}
