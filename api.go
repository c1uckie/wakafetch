package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SummaryResponse struct {
	Data []struct {
		Branches []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"branches"`
		Categories []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"categories"`
		Dependencies []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"dependencies"`
		Editors []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"editors"`
		Entities []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"entities"`
		GrandTotal struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"grand_total"`
		Languages []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"languages"`
		Machines []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"machines"`
		OperatingSystems []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"operating_systems"`
		Projects []struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Name         string  `json:"name"`
			Percent      float64 `json:"percent"`
			Seconds      int     `json:"seconds"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"projects"`
		Range struct {
			Date     string `json:"date"`
			End      string `json:"end"`
			Start    string `json:"start"`
			Text     string `json:"text"`
			Timezone string `json:"timezone"`
		} `json:"range"`
	} `json:"data"`
}

func fetchStats(apiKey, apiURL string, days int) (*SummaryResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	today := time.Now()
	todayDate := today.Format("2006-01-02")
	startDate := today.AddDate(0, 0, -days+1).Format("2006-01-02")
	requestURL := fmt.Sprintf("%s/v1/users/current/summaries?start=%s&end=%s", apiURL, startDate, todayDate)
	fmt.Println("Request URL:", requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	encodedKey := base64.StdEncoding.EncodeToString([]byte(apiKey))
	req.Header.Set("Authorization", "Basic "+encodedKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request (%s): %w", requestURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status: %s", resp.Status)
	}

	var apiResponse SummaryResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return &apiResponse, nil
}
