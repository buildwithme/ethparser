package main

import (
	"os"
	"strings"

	"github.com/buildwithme/ethparser/internal/parser"
	"github.com/buildwithme/ethparser/pkg/constants"
	"github.com/buildwithme/ethparser/pkg/logger"
)

// subscribeEnvAddresses pulls addresses from ADDRESSES in .env
func subscribeEnvAddresses(parser parser.Parser, log *logger.Logger) {
	addrs := os.Getenv(constants.ENV_ADDRESSES)
	if addrs == "" {
		return
	}

	for _, a := range strings.Split(addrs, ",") {
		a = strings.TrimSpace(a)
		if a != "" {
			parser.Subscribe(a)
			log.Printf("[INFO] Subscribed env address: %s", a)
		}
	}
}
