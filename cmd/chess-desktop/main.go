package main

import (
	"log"

	"github.com/KaiserSin/go-chess-ai/internal/infrastructure/bootstrap"
)

// main starts the desktop chess application.
func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
