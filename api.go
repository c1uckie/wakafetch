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

type DayData struct {
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
}

type SummaryResponse struct {
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
	Data  []DayData `json:"data"`
	End   string    `json:"end"`
	Start string    `json:"start"`
}

type StatsResponse struct {
	Data struct {
		Range                     string     `json:"range"`
		Status                    string     `json:"status"`
		TotalSeconds              float64    `json:"total_seconds"`
		UserID                    string     `json:"user_id"`
		Username                  string     `json:"username"`
		DailyAverage              float64    `json:"daily_average"`
		DaysIncludingHolidays     int        `json:"days_including_holidays"`
		Start                     string     `json:"start"`
		End                       string     `json:"end"`
		HumanReadableDailyAverage string     `json:"human_readable_daily_average"`
		HumanReadableRange        string     `json:"human_readable_range"`
		HumanReadableTotal        string     `json:"human_readable_total"`
		IsCodingActivityVisible   bool       `json:"is_coding_activity_visible"`
		IsOtherUsageVisible       bool       `json:"is_other_usage_visible"`
		Branches                  []StatItem `json:"branches"`
		Categories                []StatItem `json:"categories"`
		Editors                   []StatItem `json:"editors"`
		Languages                 []StatItem `json:"languages"`
		Machines                  []StatItem `json:"machines"`
		OperatingSystems          []StatItem `json:"operating_systems"`
		Projects                  []StatItem `json:"projects"`
	} `json:"data"`
}

func fetchSummary(apiKey, apiURL string, days int) (*SummaryResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	today := time.Now()
	todayDate := today.Format("2006-01-02")
	startDate := today.AddDate(0, 0, -days+1).Format("2006-01-02")
	requestURL := fmt.Sprintf("%s/v1/users/current/summaries?start=%s&end=%s", apiURL, startDate, todayDate)
	response, err := fetchApi[SummaryResponse](apiKey, requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return response, nil
}

func fetchStats(apiKey, apiURL, rangeStr string) (*StatsResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	requestURL := fmt.Sprintf("%s/v1/users/current/stats/%s", apiURL, rangeStr)
	// fmt.Println(requestURL)
	response, err := fetchApi[StatsResponse](apiKey, requestURL)
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
