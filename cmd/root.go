package cmd

import (
	"github.com/spf13/cobra"
	"ytcw/internal/config"
)

var rootCmd = &cobra.Command{
	Use: "ytcw",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(func() { _ = config.LoadConfig })

	mainGroup := cobra.Group{ID: "main", Title: "Main Commands"}
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(daemonCmd)

	adminGroup := cobra.Group{ID: "admin", Title: "Admin Commands"}
	rootCmd.AddCommand(adminAddChannelCmd)

	rootCmd.AddGroup(&mainGroup, &adminGroup)
}
