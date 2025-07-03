package main

import (
	"log"
	"ytcw/internal/db"
	"ytcw/internal/schema"
)

func main() {
	database := db.Connect()

	err := database.AutoMigrate(&schema.Category{}, &schema.Channel{}, &schema.Video{})
	if err != nil {
		log.Fatal("Failed to auto migrate db", err)
	}

	log.Println("Migration done!")
}
