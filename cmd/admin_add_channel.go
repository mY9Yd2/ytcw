package cmd

import (
	"github.com/google/uuid"
	"github.com/mY9Yd2/ytcw/internal/content"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/fetcher"
	"github.com/mY9Yd2/ytcw/internal/logger"
	"github.com/spf13/cobra"
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
	log := logger.Pretty

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
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	categoryRepo := content.NewCategoryRepository(dbCon)
	var categoryID *uuid.UUID

	if category != "" {
		id, err := categoryRepo.SaveCategory(category)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to save category")
		}
		categoryID = &id
	}

	ytFetcher := fetcher.Fetcher{
		Logger: log,
	}

	channelInfo := ytFetcher.GetChannelInfo(channel)

	info := content.Channel{
		UploaderID: channelInfo.UploaderID,
		ChannelID:  channelInfo.ChannelID,
		Channel:    channelInfo.Channel,
	}
	info.CategoryRefer = categoryID

	channelRepo := content.NewChannelRepository(dbCon)
	if err := channelRepo.SaveChannel(&info); err != nil {
		log.Fatal().Err(err).Msg("Failed to save channel")
	}

	log.Info().Msgf("Channel %s added successfully", info.Channel)
}
