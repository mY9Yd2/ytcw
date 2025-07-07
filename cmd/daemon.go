package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/config"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/fetcher"
	"github.com/mY9Yd2/ytcw/internal/mapper"
	"github.com/mY9Yd2/ytcw/internal/repository"
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
	repo := repository.Repository{DB: dbCon}
	cfg := config.GetConfig()
	ytFetcher := fetcher.Fetcher{
		Logger: log.Logger,
	}

	log.Info().Msg("Starting fetcher daemon")

	for {
		channel, err := repo.GetStaleChannel(cfg.Fetcher.MaxLastFetchAge)
		if err != nil {
			time.Sleep(cfg.Fetcher.NoChannelRetryInterval)
			continue
		}

		videoStream := ytFetcher.FetchVideos(channel.UploaderID)

		for info := range videoStream {
			video := mapper.MapVideoInfoToVideo(info)
			video.Channel = *channel
			if err := repo.SaveVideo(&video); err != nil {
				log.Warn().Err(err).Str("DisplayID", video.DisplayID).Msg("Failed to save video")
			}
		}

		channel.LastFetch = func(t time.Time) *time.Time {
			return &t
		}(time.Now().UTC())

		if err := repo.DB.Save(channel).Error; err != nil {
			log.Warn().Err(err).Str("UploaderID", channel.UploaderID).Msg("Failed to update last_fetch")
		}

		time.Sleep(cfg.Fetcher.NoChannelRetryInterval)
	}
}
