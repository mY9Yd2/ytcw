package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ytcw",
}

func init() {
	cobra.OnInitialize(func() {
		if err := config.LoadConfig(); err != nil {
			log.Fatal().Err(err).Msg("Failed to read config")
		}
	})

	mainGroup := cobra.Group{ID: "main", Title: "Main Commands"}
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(daemonCmd)

	adminGroup := cobra.Group{ID: "admin", Title: "Admin Commands"}
	rootCmd.AddCommand(adminAddChannelCmd)
	rootCmd.AddCommand(adminDeleteChannelCmd)
	rootCmd.AddCommand(adminDisableChannelCmd)

	rootCmd.AddGroup(&mainGroup, &adminGroup)
}

func Execute() error {
	return rootCmd.Execute()
}
