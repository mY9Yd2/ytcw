package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"ytcw/internal/db"
	"ytcw/internal/schema"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Run database migrations",
	Run:     migrate,
	GroupID: "main",
}

func migrate(cmd *cobra.Command, args []string) {
	database := db.Connect()

	err := database.AutoMigrate(&schema.Category{}, &schema.Channel{}, &schema.Video{})
	if err != nil {
		log.Fatal("Failed to auto migrate db", err)
	}

	fmt.Println("Migration done!")
}
