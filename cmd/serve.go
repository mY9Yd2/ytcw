package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Run:     serve,
	GroupID: "main",
}

func serve(cmd *cobra.Command, args []string) {
	log.Println("TODO")
}
