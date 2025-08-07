package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func parseConfig() (string, string, error) {
	configPath := getConfigPath()

	file, err := os.Open(configPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var apiURL, apiKey string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "api_url") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				apiURL = strings.TrimSpace(parts[1])
			}
		} else if strings.HasPrefix(line, "api_key") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				apiKey = strings.TrimSpace(parts[1])
			}
		}

		if apiURL != "" && apiKey != "" {
			break
		}
	}

	if apiURL == "" {
		apiURL = "https://wakapi.dev/api"
	}

	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}

	if apiKey == "" {
		return "", "", fmt.Errorf("api_key not found in config")
	}

	return apiURL, apiKey, nil
}

func getConfigPath() string {
	var configFile string

	if runtime.GOOS == "windows" {
		configFile = filepath.Join(os.Getenv("USERPROFILE"), ".wakatime.cfg")
	} else {
		configFile = filepath.Join(os.Getenv("HOME"), ".wakatime.cfg")
	}

	return configFile
}
