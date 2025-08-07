package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SummaryResponse struct {
	CachedAt string `json:"cached_at"`
	Data     struct {
		Branches []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"branches"`
		Categories []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"categories"`
		Dependencies []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"dependencies"`
		Editors []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"editors"`
		Entities []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"entities"`
		GrandTotal struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"grand_total"`
		Languages []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"languages"`
		Machines []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"machines"`
		OperatingSystems []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
		} `json:"operating_systems"`
		Projects []struct {
			Digital      string `json:"digital"`
			Hours        int    `json:"hours"`
			Minutes      int    `json:"minutes"`
			Name         string `json:"name"`
			Percent      int    `json:"percent"`
			Seconds      int    `json:"seconds"`
			Text         string `json:"text"`
			TotalSeconds int    `json:"total_seconds"`
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
	today := time.Now()
	todayDate := today.Format("2006-01-02")
	startDate := today.AddDate(0, 0, -days+1).Format("2006-01-02")
	requestURL := fmt.Sprintf("%s/users/current/summaries?start=%s&end=%s", apiURL, startDate, todayDate)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+apiKey)

	client := &http.Client{Timeout: 10 * 1000}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status: %s", resp.Status)
	}

	var apiResponse SummaryResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return &apiResponse, nil
}
