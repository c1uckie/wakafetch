package main

import (
	"flag"
	"fmt"
)

func main() {
	daysFlag := flag.Int("days", 1, "Number of days to fetch data for (default: 1)")
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

	data, err := fetchStats(apiKey, apiURL, *daysFlag)
	if err != nil {
		fmt.Printf("\033[31m%v\033[0m\n", err)
		return
	}
	prettyPrint(data, *fullFlag)
}
