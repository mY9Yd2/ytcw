package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"time"
	"ytcw/internal/config"
	"ytcw/internal/db"
	"ytcw/internal/fetcher"
	"ytcw/internal/mapper"
	"ytcw/internal/repository"
)

var daemonCmd = &cobra.Command{
	Use:     "daemon",
	Short:   "Start the fetcher daemon",
	Run:     daemon,
	GroupID: "main",
}

func daemon(cmd *cobra.Command, args []string) {
	repo := repository.Repository{DB: db.Connect()}
	cfg := config.LoadConfig()

	log.Info().Msg("Starting fetcher daemon")

	for {
		channel, err := repo.GetStaleChannel(cfg.Ytcwd.MaxLastFetchAge)
		if err != nil {
			time.Sleep(cfg.Ytcwd.NoChannelRetryInterval)
			continue
		}

		videoStream := fetcher.FetchVideos(channel.UploaderID)

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

		time.Sleep(cfg.Ytcwd.NoChannelRetryInterval)
	}
}
