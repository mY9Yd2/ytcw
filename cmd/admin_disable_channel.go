package cmd

import (
	"fmt"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/logger"
	"github.com/mY9Yd2/ytcw/internal/repository"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var adminDisableChannelCmd = &cobra.Command{
	Use:     "disable-channel",
	Short:   "Disable a channel",
	Run:     disableChannel,
	Args:    disableChannelArgsValidator,
	GroupID: "admin",
}

func init() {
	adminDisableChannelCmd.Flags().StringP("id", "i", "", "Channel ID or @handle (required)")
	_ = adminDisableChannelCmd.MarkFlagRequired("id")

	// ~80 years. This app probably won't live that long,
	// but if someone wants to disable a channel "forever", it's a reasonable default.
	yearsToDisable := 80 * 365 * 24 * time.Hour

	adminDisableChannelCmd.Flags().DurationP("duration", "d", yearsToDisable, "Disable a channel for a given duration")
}

func disableChannelArgsValidator(cmd *cobra.Command, args []string) error {
	disableDuration, err := cmd.Flags().GetDuration("duration")
	if err != nil {
		return fmt.Errorf("failed to get 'duration' flag")
	}

	if disableDuration <= 0 {
		return fmt.Errorf("duration cannot be negative or zero")
	}

	return nil
}

func disableChannel(cmd *cobra.Command, args []string) {
	log := logger.Pretty

	channel, err := cmd.Flags().GetString("id")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'id' flag")
	}

	disableDuration, _ := cmd.Flags().GetDuration("duration")

	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	channelRepo := repository.NewChannelRepository(dbCon)

	now := time.Now().UTC()
	disabledUntil := now.Add(disableDuration)

	if strings.HasPrefix(channel, "@") {
		if err := channelRepo.DisableChannelByUploaderID(channel, now, disabledUntil); err != nil {
			log.Fatal().Err(err).Msg("Failed to disable channel")
		}
	} else {
		if err := channelRepo.DisableChannelByChannelID(channel, now, disabledUntil); err != nil {
			log.Fatal().Err(err).Msg("Failed to disable channel")
		}
	}

	log.Info().Msgf("Channel %s disabled until %s", channel, disabledUntil.Format(time.RFC3339))
}
