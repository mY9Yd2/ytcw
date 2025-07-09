package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/config"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/fetcher"
	"github.com/mY9Yd2/ytcw/internal/repository"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"time"
)

var daemonCmd = &cobra.Command{
	Use:     "daemon",
	Short:   "Start the fetcher daemon",
	Run:     daemon,
	GroupID: "main",
}

func daemon(cmd *cobra.Command, args []string) {
	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	channelRepo := repository.NewChannelRepository(dbCon)
	videoRepo := repository.NewVideoRepository(dbCon)
	cfg := config.GetConfig()
	ytFetcher := fetcher.Fetcher{
		Logger: log.Logger,
	}

	log.Info().Msg("Starting fetcher daemon")

	for {
		channel, err := channelRepo.GetStaleChannel(cfg.Fetcher.MaxLastFetchAge)
		if err != nil {
			time.Sleep(cfg.Fetcher.NoChannelRetryInterval)
			continue
		}

		videoStream := ytFetcher.FetchVideos(channel.UploaderID)

		for info := range videoStream {
			video := schema.Video{
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
				}}
			video.Channel = *channel
			if err := videoRepo.SaveVideo(&video); err != nil {
				log.Warn().Err(err).Str("DisplayID", video.DisplayID).Msg("Failed to save video")
			}
		}

		err = channelRepo.UpdateChannelLastFetch(channel.ID, time.Now().UTC())
		if err != nil {
			log.Warn().Err(err).Msg("Failed to update channel last fetch")
		}

		time.Sleep(cfg.Fetcher.NoChannelRetryInterval)
	}
}
