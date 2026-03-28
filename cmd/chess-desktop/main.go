package main

import (
	"log"

	"github.com/KaiserSin/go-chess-ai/internal/infrastructure/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
