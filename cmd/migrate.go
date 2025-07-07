package cmd

import (
	"fmt"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/schema"
	"github.com/spf13/cobra"
	"log"
)

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Short:   "Run database migrations",
	Run:     migrate,
	GroupID: "main",
}

func migrate(cmd *cobra.Command, args []string) {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = database.AutoMigrate(&schema.Category{}, &schema.Channel{}, &schema.Video{})
	if err != nil {
		log.Fatal("Failed to auto migrate db", err)
	}

	fmt.Println("Migration done!")
}
