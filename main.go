package main

import (
	"github.com/mY9Yd2/ytcw/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
