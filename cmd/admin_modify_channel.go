package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/logger"
	"github.com/mY9Yd2/ytcw/internal/repository"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"github.com/spf13/cobra"
	"strings"
)

var adminModifyChannelCmd = &cobra.Command{
	Use:     "modify-channel",
	Short:   "Modify a channel",
	Run:     modifyChannel,
	GroupID: "admin",
}

func init() {
	adminModifyChannelCmd.Flags().StringP("id", "i", "", "Channel ID or @handle (required)")
	_ = adminModifyChannelCmd.MarkFlagRequired("id")

	adminModifyChannelCmd.Flags().StringP("category", "c", "", "Set the channel's category")
	adminModifyChannelCmd.Flags().Bool("unset-category", false, "Unset the channel's category")
}

func modifyChannel(cmd *cobra.Command, args []string) {
	log := logger.Pretty

	channel, err := cmd.Flags().GetString("id")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'id' flag")
	}

	category, err := cmd.Flags().GetString("category")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'category' flag")
	}

	unsetCategory, err := cmd.Flags().GetBool("unset-category")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'unset-category' flag")
	}

	if category != "" && unsetCategory {
		log.Fatal().Msg("Flags --category and --unset-category cannot be used together")
	}

	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	channelRepo := repository.NewChannelRepository(dbCon)

	var existingChannel *schema.Channel
	if strings.HasPrefix(channel, "@") {
		existingChannel, err = channelRepo.GetChannelByUploaderID(channel)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to find channel")
		}
	} else {
		existingChannel, err = channelRepo.GetChannelByChannelID(channel)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to find channel")
		}
	}

	if category != "" {
		categoryRepo := repository.NewCategoryRepository(dbCon)
		id, err := categoryRepo.SaveCategory(category)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to save category")
		}
		existingChannel.CategoryRefer = &id
	} else if unsetCategory {
		existingChannel.CategoryRefer = nil
	}

	if err := channelRepo.SaveChannel(existingChannel); err != nil {
		log.Fatal().Err(err).Msg("Failed to save updated channel")
	}

	log.Info().Msg("Channel updated successfully")
}
