package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sahaj-b/wakafetch/ui"
)

func main() {
	// daysFlag := flag.Int("days", 1, "Number of days to fetch data for")
	rangeFlag := flag.String("range", "7d", "Range of data to fetch (today/7d/30d/6m/1y/all)")
	apiKeyFlag := flag.String("api-key", "", "Your WakaTime/Wakapi API key (overrides config)")
	fullFlag := flag.Bool("full", false, "Display full statistics")
	noColorFlag := flag.Bool("no-colors", false, "Disable coloRed output")
	helpFlag := flag.Bool("help", false, "Display help information")
	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage: wakafetch [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	if !*noColorFlag && colorsShouldBeEnabled() {
		ui.EnableColors()
	}

	apiURL, apiKey, err := parseConfig()
	if err != nil {
		ui.Errorln("%v", err)
		os.Exit(1)
		return
	}

	if *apiKeyFlag != "" {
		apiKey = *apiKeyFlag
	}

	rangeStrMap := map[string]string{
		"today": "today",
		"7d":    "last_7_days",
		"30d":   "last_30_days",
		"6m":    "last_6_months",
		"1y":    "last_year",
		"all":   "all_time",
	}
	rangeStr, exists := rangeStrMap[*rangeFlag]
	if !exists {
		ui.Errorln("Invalid range: '%s', must be one of %stoday, 7d, 30d, 6m, 1y, all", *rangeFlag, ui.Clr.Green)
		os.Exit(1)
		return
	}

	if rangeStr == "all_time" {
		data, err := fetchStats(apiKey, apiURL, rangeStr)
		if err != nil {
			ui.Errorln("%v", err)
			os.Exit(1)
			return
		}
		ui.DisplayStats(data, *fullFlag, rangeStr)
		return
	}

	days := map[string]int{
		"today":         1,
		"last_7_days":   7,
		"last_30_days":  30,
		"last_6_months": 183,
		"last_year":     365,
	}[rangeStr]

	data, err := fetchSummary(apiKey, apiURL, days)
	if err != nil {
		ui.Errorln("%v", err)
		os.Exit(1)
		return
	}

	if len(data.Data) == 0 {
		ui.Warnln("No data available for the selected period.")
		return
	}
	ui.DisplaySummary(data, *fullFlag, rangeStr)
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
