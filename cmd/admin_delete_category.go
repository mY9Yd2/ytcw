package cmd

import (
	"github.com/mY9Yd2/ytcw/internal/content"
	"github.com/mY9Yd2/ytcw/internal/db"
	"github.com/mY9Yd2/ytcw/internal/logger"
	"github.com/spf13/cobra"
)

var adminDeleteCategoryCmd = &cobra.Command{
	Use:     "delete-category",
	Short:   "Delete an empty category",
	Run:     deleteCategory,
	GroupID: "admin",
}

func init() {
	adminDeleteCategoryCmd.Flags().StringP("category", "c", "", "Category name (required)")
	_ = adminDeleteCategoryCmd.MarkFlagRequired("category")
}

func deleteCategory(cmd *cobra.Command, args []string) {
	log := logger.Pretty

	category, err := cmd.Flags().GetString("category")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get 'category' flag")
	}

	dbCon, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	categoryRepo := content.NewCategoryRepository(dbCon)

	if err = categoryRepo.DeleteCategory(category); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete category")
	}

	log.Info().Msgf("Category %s deleted successfully", category)
}
