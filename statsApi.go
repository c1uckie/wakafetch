package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type StatsResponse struct {
	Data struct {
		Range                     string     `json:"range"`
		Status                    string     `json:"status"`
		TotalSeconds              int        `json:"total_seconds"`
		UserID                    string     `json:"user_id"`
		Username                  string     `json:"username"`
		DailyAverage              int        `json:"daily_average"`
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

func fetchStats(apiKey, apiURL, rangeStr string) (*StatsResponse, error) {
	apiURL = strings.TrimSuffix(apiURL, "/")
	requestURL := fmt.Sprintf("%s/v1/users/current/stats/%s", apiURL, rangeStr)
	fmt.Println(requestURL)
	response, err := fetchApi[StatsResponse](apiKey, requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	return response, nil
}

