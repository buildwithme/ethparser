package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadDotEnv parses a .env file and sets environment variables in the process.
func LoadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
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
