package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"ytcw/internal/db"
	"ytcw/internal/fetcher"
	"ytcw/internal/mapper"
	"ytcw/internal/repository"
)

var channel string
var category string

var adminAddChannelCmd = &cobra.Command{
	Use:     "add-channel",
	Short:   "Add a new channel",
	Run:     add,
	GroupID: "admin",
}

func init() {
	adminAddChannelCmd.Flags().StringVarP(&channel, "id", "i", "", "Channel ID or @handle (required)")
	_ = adminAddChannelCmd.MarkFlagRequired("id")

	adminAddChannelCmd.Flags().StringVarP(&category, "category", "c", "", "Category")

	rootCmd.AddCommand(adminAddChannelCmd)
}

func add(cmd *cobra.Command, args []string) {
	repo := repository.Repository{DB: db.Connect()}
	var categoryID *uint

	if category != "" {
		id, err := repo.SaveCategory(category)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to save category")
		}
		categoryID = &id
	}

	info := mapper.MapChannelInfoToChannel(fetcher.GetChannelInfo(channel))
	info.CategoryRefer = categoryID

	if err := repo.SaveChannel(&info); err != nil {
		log.Fatal().Err(err).Msg("Failed to save channel")
	}
}
