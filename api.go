package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sahaj-b/wakafetch/types"
)

func fetchSummary(apiKey, apiURL string, days int) (*types.SummaryResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	today := time.Now()
	todayDate := today.Format("2006-01-02")
	startDate := today.AddDate(0, 0, -days+1).Format("2006-01-02")
	requestURL := fmt.Sprintf("%s/v1/users/current/summaries?start=%s&end=%s", apiURL, startDate, todayDate)
	response, err := fetchApi[types.SummaryResponse](apiKey, requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return response, nil
}

func fetchStats(apiKey, apiURL, rangeStr string) (*types.StatsResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	requestURL := fmt.Sprintf("%s/v1/users/current/stats/%s", apiURL, rangeStr)
	// fmt.Println(requestURL)
	response, err := fetchApi[types.StatsResponse](apiKey, requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return response, nil
}

func fetchApi[T any](apiKey, requestURL string) (*T, error) {
	const timeout = 10 * time.Second
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	encodedKey := base64.StdEncoding.EncodeToString([]byte(apiKey))
	req.Header.Set("Authorization", "Basic "+encodedKey)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request (%s): %w", requestURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status: %s", resp.Status)
	}

	var apiResponse T
	if err = json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return &apiResponse, nil
}
