package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"ytcw/internal/db"
	"ytcw/internal/fetcher"
	"ytcw/internal/logger"
	"ytcw/internal/mapper"
	"ytcw/internal/repository"
)

var adminAddChannelCmd = &cobra.Command{
	Use:     "add-channel",
	Short:   "Add a new channel",
	Run:     addChannel,
	GroupID: "admin",
}

func init() {
	adminAddChannelCmd.Flags().StringP("id", "i", "", "Channel ID or @handle (required)")
	_ = adminAddChannelCmd.MarkFlagRequired("id")

	adminAddChannelCmd.Flags().StringP("category", "c", "", "Category")
}

func addChannel(cmd *cobra.Command, args []string) {
	channel, err := cmd.Flags().GetString("id")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'id' flag")
	}

	category, err := cmd.Flags().GetString("category")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'category' flag")
	}

	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	repo := repository.Repository{DB: dbCon}
	var categoryID *uint

	if category != "" {
		id, err := repo.SaveCategory(category)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to save category")
		}
		categoryID = &id
	}

	ytFetcher := fetcher.Fetcher{
		Logger: logger.JSON,
	}
	info := mapper.MapChannelInfoToChannel(ytFetcher.GetChannelInfo(channel))
	info.CategoryRefer = categoryID

	if err := repo.SaveChannel(&info); err != nil {
		log.Fatal().Err(err).Msg("Failed to save channel")
	}
}
