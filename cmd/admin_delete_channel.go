package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/repository"
	"github.com/rs/zerolog/log"
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
	channel, err := cmd.Flags().GetString("id")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'id' flag")
	}

	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	repo := repository.Repository{DB: dbCon}

	if strings.HasPrefix(channel, "@") {
		if err := repo.SoftDeleteChannelByUploaderID(channel); err != nil {
			log.Fatal().Err(err).Msg("Failed to delete channel")
		}
	} else {
		if err := repo.SoftDeleteChannelByChannelID(channel); err != nil {
			log.Fatal().Err(err).Msg("Failed to delete channel")
		}
	}
}
