package main

import (
	"flag"
	"fmt"
)

func main() {
	// daysFlag := flag.Int("days", 1, "Number of days to fetch data for")
	rangeFlag := flag.String("range", "last_7_days", "Range of data to fetch (last_7_days, last_30_days, last_6_months, last_year, all_time)")
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

	data, err := fetchStats(apiKey, apiURL, *rangeFlag)
	if err != nil {
		fmt.Printf("\033[31m%v\033[0m\n", err)
		return
	}
	displayStats(data, *fullFlag, *rangeFlag)

	// data, err := fetchSummary(apiKey, apiURL, *daysFlag)
	// if err != nil {
	// 	fmt.Printf("\033[31m%v\033[0m\n", err)
	// 	return
	// }
	// displaySummary(data, *fullFlag, *daysFlag)
}
