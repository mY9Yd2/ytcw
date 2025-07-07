package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/logger"
	"github.com/mY9Yd2/ytcw/internal/repository"
	"github.com/spf13/cobra"
	"strings"
)

var adminDeleteChannelCmd = &cobra.Command{
	Use:     "delete-channel",
	Short:   "Delete a channel",
	Run:     deleteChannel,
	GroupID: "admin",
}

func init() {
	adminDeleteChannelCmd.Flags().StringP("id", "i", "", "Channel ID or @handle (required)")
	_ = adminDeleteChannelCmd.MarkFlagRequired("id")
}

func deleteChannel(cmd *cobra.Command, args []string) {
	log := logger.Pretty

	channel, err := cmd.Flags().GetString("id")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'id' flag")
	}

	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	channelRepo := repository.NewChannelRepository(dbCon)

	if strings.HasPrefix(channel, "@") {
		if err := channelRepo.SoftDeleteChannelByUploaderID(channel); err != nil {
			log.Fatal().Err(err).Msg("Failed to delete channel")
		}
	} else {
		if err := channelRepo.SoftDeleteChannelByChannelID(channel); err != nil {
			log.Fatal().Err(err).Msg("Failed to delete channel")
		}
	}
}
