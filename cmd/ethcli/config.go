package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/buildwithme/ethparser/pkg/constants"
)

// ConfigFlags holds possible command-line overrides.
type ConfigFlags struct {
	EnvFile     string
	Addresses   string
	RPCEndpoint string
	Concurrency int
	ChunkSize   int
	MaxRetries  int
	StartBlock  int
	EndBlock    int
}

// ParseFlags parses CLI flags and returns them in ConfigFlags.
// If no flag is given, it retains zero/empty values (meaning "no override").
func ParseFlags() ConfigFlags {
	var cf ConfigFlags

	flag.StringVar(&cf.EnvFile, "env", ".env", "Override the .env file path (default: .env).")
	flag.StringVar(&cf.Addresses, "addresses", "", "Override the ADDRESSES env var (default from .env).")
	flag.StringVar(&cf.RPCEndpoint, "rpc", "", "Override the RPC_ENDPOINT env var (default from .env).")
	flag.IntVar(&cf.Concurrency, "concurrency", 0, "Override the CONCURRENCY env var (default from .env).")
	flag.IntVar(&cf.ChunkSize, "chunk-size", 0, "Override the CHUNK_SIZE env var (default from .env).")
	flag.IntVar(&cf.MaxRetries, "max-retries", 0, "Override the MAX_RETRIES env var (default from .env).")
	flag.IntVar(&cf.StartBlock, "start", 0, "Override the DEFAULT_START_BLOCK env var (default from .env).")
	flag.IntVar(&cf.EndBlock, "end", 0, "Override the DEFAULT_END_BLOCK env var (default from .env).")

	// Customize usage help if desired:
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	return cf
}

// ApplyEnvFile sets the ENV_FILE_PATH environment variable if the flag was provided.
func (cf ConfigFlags) ApplyEnvFile() {
	if cf.EnvFile != "" {
		os.Setenv(constants.ENV_FILE_PATH, cf.EnvFile)
	}
}

// ApplyConfig overrides environment variables if corresponding flags were provided.
func (cf ConfigFlags) ApplyConfig() {
	if cf.Addresses != "" {
		os.Setenv(constants.ENV_ADDRESSES, cf.Addresses)
	}

	if cf.RPCEndpoint != "" {
		os.Setenv(constants.ENV_RPC_ENDPOINT, cf.RPCEndpoint)
	}

	if cf.Concurrency > 0 {
		os.Setenv(constants.ENV_CONCURRENCY, strconv.Itoa(cf.Concurrency))
	}

	if cf.ChunkSize > 0 {
		os.Setenv(constants.ENV_CHUNK_SIZE, strconv.Itoa(cf.ChunkSize))
	}

	if cf.MaxRetries > 0 {
		os.Setenv(constants.ENV_MAX_RETRIES, strconv.Itoa(cf.MaxRetries))
	}

	if cf.StartBlock > 0 {
		os.Setenv(constants.ENV_DEFAULT_START_BLOCK, strconv.Itoa(cf.StartBlock))
	}

	if cf.EndBlock > 0 {
		os.Setenv(constants.ENV_DEFAULT_END_BLOCK, strconv.Itoa(cf.EndBlock))
	}
}
