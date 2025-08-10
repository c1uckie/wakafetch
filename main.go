package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	red      = "\x1b[31m"
	yellow   = "\x1b[33m"
	boldBlue = "\x1b[34;1m"
	green    = "\x1b[32m"
	gray     = "\x1b[90m"
	bold     = "\x1b[1m"
	reset    = "\x1b[0m"
)

func main() {
	// daysFlag := flag.Int("days", 1, "Number of days to fetch data for")
	rangeFlag := flag.String("range", "7d", "Range of data to fetch (today/7d/30d/6m/1y/all)")
	apiKeyFlag := flag.String("api-key", "", "Your WakaTime/Wakapi API key (overrides config)")
	fullFlag := flag.Bool("full", false, "Display full statistics")
	noColorFlag := flag.Bool("no-colors", false, "Disable colored output")
	helpFlag := flag.Bool("help", false, "Display help information")
	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage: wakafetch [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	if *noColorFlag || !colorsShouldBeEnabled() {
		disableColors()
	}

	apiURL, apiKey, err := parseConfig()
	if err != nil {
		errorln("%v", err)
		os.Exit(1)
		return
	}

	if *apiKeyFlag != "" {
		apiKey = *apiKeyFlag
	}

	rangeStr := "last_7_days"
	switch *rangeFlag {
	case "7d":
		rangeStr = "last_7_days"
	case "30d":
		rangeStr = "last_30_days"
	case "6m":
		rangeStr = "last_6_months"
	case "1y":
		rangeStr = "last_year"
	case "all":
		rangeStr = "all_time"
	case "today":
		rangeStr = "today"
	default:
		errorln("Invalid range: '%s', must be one of %stoday, 7d, 30d, 6m, 1y, all",
			*rangeFlag, green)
		os.Exit(1)
		return
	}
	if rangeStr == "today" {
		data, err := fetchSummary(apiKey, apiURL, 1)
		if err != nil {
			errorln("%v", err)
			os.Exit(1)
			return
		}
		displaySummary(data, *fullFlag, 1)
	} else {
		data, err := fetchStats(apiKey, apiURL, rangeStr)
		if err != nil {
			errorln("%v", err)
			os.Exit(1)
			return
		}
		displayStats(data, *fullFlag, rangeStr)
	}
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

func disableColors() {
	red, yellow, boldBlue, green, gray, bold, reset = "", "", "", "", "", "", ""
}
