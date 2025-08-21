package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sahaj-b/wakafetch/ui"
)

// okay so there are 2 types of responses/endpoints:
// /summary response:
// - gives summary of EACH day, so more granular data
// - have to aggregate data manually for viewing stats
// - supports custom date ranges
// - so, its called when --days(custom range) or --daily/heatmap(granular daily breakdown) is used

// /stats response:
// - gives summary of the ENTIRE range in a single response
// - no need for aggregation, efficient af
// - doesn't support custom date ranges(only rangeStr)
// - so, its the default unless --days or --daily/heatmap is used

func main() {
	config := parseFlags()
	apiURL, apiKey := loadAPIConfig(config)

	if shouldUseSummaryAPI(config) {
		handleSummaryFlow(config, apiKey, apiURL)
	} else {
		handleStatsFlow(config, apiKey, apiURL)
	}
}

func loadAPIConfig(config Config) (string, string) {
	apiURL, apiKey, err := parseConfig()
	if err != nil {
		ui.Errorln(err.Error())
	}

	if *config.apiKeyFlag != "" {
		apiKey = *config.apiKeyFlag
	}

	return apiURL, apiKey
}

func shouldUseSummaryAPI(config Config) bool {
	return *config.daysFlag != 0 || *config.dailyFlag || *config.heatmapFlag
}

func handleStatsFlow(config Config, apiKey, apiURL string) {
	rangeStr := getRangeStr(*config.rangeFlag)

	data, err := fetchStats(apiKey, apiURL, rangeStr)
	if err != nil {
		ui.Errorln(err.Error())
	}

	if *config.jsonFlag {
		outputJSON(data)
		return
	}

	ui.DisplayStats(data, *config.fullFlag, rangeStr)
}

func handleSummaryFlow(config Config, apiKey, apiURL string) {
	rangeStr := getRangeStr(*config.rangeFlag)
	days := *config.daysFlag
	validRange := true
	if days == 0 {
		days, validRange = map[string]int{
			"today":         1,
			"last_7_days":   7,
			"last_30_days":  30,
			"last_6_months": 183,
			"last_year":     365,
		}[rangeStr]
	}

	if !validRange {
		ui.Errorln("This range isn't supported with `--daily` or `--heatmap` flags. Use `--days` instead")
		return
	}

	data, err := fetchSummary(apiKey, apiURL, days)
	if err != nil {
		ui.Errorln(err.Error())
	}

	if *config.jsonFlag {
		outputJSON(data)
		return
	}

	var heading string
	if *config.daysFlag != 0 {
		if days == 1 {
			heading = "Today"
		} else {
			heading = fmt.Sprintf("Last %d days", days)
		}
	} else {
		headingMap := map[string]string{
			"today":         "Today",
			"last_7_days":   "Last 7 days",
			"last_30_days":  "Last 30 days",
			"last_6_months": "Last 6 months",
			"last_year":     "Last year",
			"all_time":      "All time",
		}
		heading = headingMap[rangeStr]
	}

	if *config.dailyFlag {
		ui.DisplayBreakdown(data.Data, heading)
		return
	}

	if *config.heatmapFlag {
		ui.DisplayHeatmap(data.Data, heading)
		return
	}

	ui.DisplaySummary(data, *config.fullFlag, heading)
}

func getRangeStr(rangeFlag string) string {
	rangeStrMap := map[string]string{
		"today": "today",
		"7d":    "last_7_days",
		"30d":   "last_30_days",
		"6m":    "last_6_months",
		"1y":    "last_year",
		"all":   "all_time",
	}
	rangeStr, exists := rangeStrMap[rangeFlag]
	if !exists {
		ui.Errorln("Invalid range: '%s', must be one of %stoday, 7d, 30d, 6m, 1y, all", rangeFlag, ui.Clr.Green)
	}
	return rangeStr
}

func colorsShouldBeEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}

	// tty check
	file, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (file.Mode() & os.ModeCharDevice) != 0
}

func outputJSON(data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		ui.Errorln("Failed to marshal JSON: %s", err.Error())
	}
	fmt.Println(string(jsonData))
}
