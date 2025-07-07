package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/logger"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Run database migrations",
	Run:     migrate,
	GroupID: "main",
}

func migrate(cmd *cobra.Command, args []string) {
	log := logger.Pretty

	database, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	err = database.AutoMigrate(&schema.Category{}, &schema.Channel{}, &schema.Video{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to auto migrate the database")
	}

	log.Info().Msg("Database migrations complete")
}
