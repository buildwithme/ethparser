package main

import (
	"log"

	"github.com/buildwithme/ethparser/pkg/env"
)

func main() {
	if err := env.LoadDotEnv(".env"); err != nil {
		log.Fatalf("[WARN] .env not loaded: %v", err)
	}
}
