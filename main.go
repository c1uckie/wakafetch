package main

import (
	"flag"
	"fmt"
)

func main() {
	// daysFlag := flag.Int("days", 1, "Number of days to fetch data for")
	rangeFlag := flag.String("range", "7d", "Range of data to fetch (7d/30d/6m/1y/all)")
	apiKeyFlag := flag.String("api-key", "", "Your WakaTime/Wakapi API key (overrides config)")
	helpFlag := flag.Bool("help", false, "Display help information")
	fullFlag := flag.Bool("full", false, "Display full statistics")
	flag.Parse()

	if *helpFlag {
		fmt.Println("Usage: wakafetch [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	apiURL, apiKey, err := parseConfig()

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
	default:
		fmt.Printf(red+"Invalid range: '%s', must be one of %s7d, 30d, 6m, 1y, all"+reset+"\n",
			*rangeFlag, green)
		return
	}
	data, err := fetchStats(apiKey, apiURL, rangeStr)
	if err != nil {
		fmt.Printf("\033[31m%v\033[0m\n", err)
		return
	}
	displayStats(data, *fullFlag, rangeStr)

	// data, err := fetchSummary(apiKey, apiURL, *daysFlag)
	// if err != nil {
	// 	fmt.Printf("\033[31m%v\033[0m\n", err)
	// 	return
	// }
	// displaySummary(data, *fullFlag, *daysFlag)
}
