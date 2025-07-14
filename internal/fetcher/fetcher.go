package fetcher

import (
	videotype "github.com/mY9Yd2/ytcw/internal/model"
	model "github.com/mY9Yd2/ytcw/internal/model/fetcher"
	"github.com/rs/zerolog"
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

func (f *Fetcher) FetchRegularVideos(channel string) <-chan model.VideoInfo {
	rawOut := make(chan model.VideoInfo)
	processedOut := make(chan model.VideoInfo)

	go func() {
		defer close(rawOut)

		stop := make(chan struct{})

		fetchOpts := fetchOptions{
			checkCutoff:  true,
			addThumbnail: true,
		}

		yt := newYtDlp().
			SetVideosURL(channel).
			SetFetchOpts(fetchOpts)
		if err := yt.fetch(f.Logger, rawOut, stop); err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching videos")
		}
	}()

	go func() {
		defer close(processedOut)
		for video := range rawOut {
			video.VideoType = videotype.VideoTypeRegular
			processedOut <- video
		}
	}()

	return processedOut
}

func (f *Fetcher) FetchShorts(channel string) <-chan model.VideoInfo {
	rawOut := make(chan model.VideoInfo)
	processedOut := make(chan model.VideoInfo)

	go func() {
		defer close(rawOut)

		stop := make(chan struct{})

		fetchOpts := fetchOptions{
			checkCutoff:  true,
			addThumbnail: true,
		}

		yt := newYtDlp().
			SetShortsURL(channel).
			SetFetchOpts(fetchOpts)
		if err := yt.fetch(f.Logger, rawOut, stop); err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching shorts")
		}
	}()

	go func() {
		defer close(processedOut)
		for video := range rawOut {
			video.VideoType = videotype.VideoTypeShort
			processedOut <- video
		}
	}()

	return processedOut
}
