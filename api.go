package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type StatItem struct {
	Digital      string  `json:"digital"`
	Hours        int     `json:"hours"`
	Minutes      int     `json:"minutes"`
	Name         string  `json:"name"`
	Percent      float64 `json:"percent"`
	Seconds      int     `json:"seconds"`
	Text         string  `json:"text"`
	TotalSeconds float64 `json:"total_seconds"`
}

type SummaryResponse struct {
	Data []struct {
		Branches         []StatItem `json:"branches"`
		Categories       []StatItem `json:"categories"`
		Dependencies     []StatItem `json:"dependencies"`
		Editors          []StatItem `json:"editors"`
		Entities         []StatItem `json:"entities"`
		Languages        []StatItem `json:"languages"`
		Machines         []StatItem `json:"machines"`
		OperatingSystems []StatItem `json:"operating_systems"`
		Projects         []StatItem `json:"projects"`
		GrandTotal       struct {
			Digital      string  `json:"digital"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Text         string  `json:"text"`
			TotalSeconds float64 `json:"total_seconds"`
		} `json:"grand_total"`
		Range struct {
			Date     string `json:"date"`
			End      string `json:"end"`
			Start    string `json:"start"`
			Text     string `json:"text"`
			Timezone string `json:"timezone"`
		} `json:"range"`
	} `json:"data"`
	End             string `json:"end"`
	Start           string `json:"start"`
	CumulativeTotal struct {
		Digital string `json:"digital"`
		Seconds int    `json:"seconds"`
		Text    string `json:"text"`
	} `json:"cumulative_total"`
	DailyAverage struct {
		DaysIncludingHolidays int    `json:"days_including_holidays"`
		DaysMinusHolidays     int    `json:"days_minus_holidays"`
		Holidays              int    `json:"holidays"`
		Seconds               int    `json:"seconds"`
		Text                  string `json:"text"`
	} `json:"daily_average"`
}

func fetchStats(apiKey, apiURL string, days int) (*SummaryResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	today := time.Now()
	todayDate := today.Format("2006-01-02")
	startDate := today.AddDate(0, 0, -days+1).Format("2006-01-02")
	requestURL := fmt.Sprintf("%s/v1/users/current/summaries?start=%s&end=%s", apiURL, startDate, todayDate)
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
