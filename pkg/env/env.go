package env

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/buildwithme/ethparser/pkg/constants"
)

type EnvVars struct {
}

func NewEnvVars() *EnvVars {
	return &EnvVars{}
}

// LoadDotEnv parses a .env file and sets environment variables in the process.
func LoadDotEnv() error {
	path := GetEnvString(constants.ENV_FILE_PATH, ".env")

	log.Printf("Open env filepath: %s", path)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("[error]: Failed to load path: %s", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		os.Setenv(key, val)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env: %w", err)
	}

	return nil
}

func GetEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return i
}

func GetEnvString(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	return val
}
