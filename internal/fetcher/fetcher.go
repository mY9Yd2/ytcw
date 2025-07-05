package fetcher

import (
	"github.com/rs/zerolog"
	"ytcw/internal/model"
)

type Fetcher struct {
	Logger zerolog.Logger
	ytdlp  ytdlp
}

func (f *Fetcher) GetChannelInfo(channel string) model.ChannelInfo {
	out := make(chan model.VideoInfo)
	stop := make(chan struct{})

	go func() {
		defer close(out)

		yt := newYtDlp().
			AddArg("-I:10").
			SetChannelURL(channel).
			SetFetchOpts(fetchOptions{
				checkCutoff:  false,
				addThumbnail: false,
			})

		if err := yt.fetch(f.Logger, out, stop); err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching channel info")
		}
	}()

	info, ok := <-out
	if !ok {
		f.Logger.Fatal().Msg("No video info found")
	}
	close(stop)

	return model.ChannelInfo{
		UploaderID: info.UploaderID,
		ChannelID:  info.ChannelID,
		Channel:    info.Channel,
	}
}

func (f *Fetcher) FetchVideos(channel string) <-chan model.VideoInfo {
	out := make(chan model.VideoInfo)

	go func() {
		defer close(out)

		stop := make(chan struct{})

		fetchOpts := fetchOptions{
			checkCutoff:  true,
			addThumbnail: true,
		}

		yt := newYtDlp().
			SetShortsURL(channel).
			SetFetchOpts(fetchOpts)
		if err := yt.fetch(f.Logger, out, stop); err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching shorts")
		}

		stop = make(chan struct{}) // Reset

		yt.SetVideosURL(channel)
		if err := yt.fetch(f.Logger, out, stop); err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching videos")
		}
	}()

	return out
}
