package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sahaj-b/wakafetch/ui"
)

type Config struct {
	rangeFlag   *string
	apiKeyFlag  *string
	fullFlag    *bool
	daysFlag    *int
	dailyFlag   *bool
	noColorFlag *bool
	helpFlag    *bool
}

// okay so there are 2 types of responses/endpoints:
// /summary response:
// - gives summary of EACH day, so more granular data
// - have to aggregate data manually for viewing stats
// - supports custom date ranges
// - so, its called when --days(custom range) or --daily(granular daily breakdown) is used

// /stats response:
// - gives summary of the ENTIRE range in a single response
// - no need for aggregation, efficient af
// - doesn't support custom date ranges(only rangeStr)
// - so, its the default unless --days or --daily is used

func main() {
	config := parseFlags()
	apiURL, apiKey := loadAPIConfig(config)

	if shouldUseSummaryAPI(config) {
		handleSummaryFlow(config, apiKey, apiURL)
	} else {
		handleStatsFlow(config, apiKey, apiURL)
	}
}

type flagInfo struct {
	longName    string
	shortName   string
	defaultVal  any
	description string
	flagType    string
}

var flagDefs = []flagInfo{
	{"range", "r", "7d", "Range of data to fetch (today/7d/30d/6m/1y/all) (default: 7d)", "string"},
	{"api-key", "k", "", "Your WakaTime/Wakapi API key (overrides config)", "string"},
	{"full", "f", false, "Display full statistics", ""},
	{"days", "d", 0, "Number of days to fetch data for (overrides --range)", "int"},
	{"daily", "D", false, "Display daily breakdown", ""},
	{"no-colors", "n", false, "Disable colored output", ""},
	{"help", "h", false, "Display help information", ""},
}

func parseFlags() Config {
	config := Config{
		rangeFlag:   flag.String("range", "7d", ""),
		apiKeyFlag:  flag.String("api-key", "", ""),
		fullFlag:    flag.Bool("full", false, ""),
		daysFlag:    flag.Int("days", 0, ""),
		dailyFlag:   flag.Bool("daily", false, ""),
		noColorFlag: flag.Bool("no-colors", false, ""),
		helpFlag:    flag.Bool("help", false, ""),
	}

	flag.StringVar(config.rangeFlag, "r", "7d", "")
	flag.StringVar(config.apiKeyFlag, "k", "", "")
	flag.BoolVar(config.fullFlag, "f", false, "")
	flag.IntVar(config.daysFlag, "d", 0, "")
	flag.BoolVar(config.dailyFlag, "D", false, "")
	flag.BoolVar(config.noColorFlag, "n", false, "")
	flag.BoolVar(config.helpFlag, "h", false, "")

	flag.Usage = showCustomHelp
	flag.Parse()

	if *config.noColorFlag || !colorsShouldBeEnabled() {
		ui.DisableColors()
	}

	if *config.helpFlag {
		showCustomHelp()
	}

	if *config.daysFlag < 0 {
		log.Fatal("Invalid value for --days: must be a positive integer")
	}

	return config
}

func showCustomHelp() {
	fmt.Println(ui.Clr.Bold + "Usage:" + ui.Clr.Reset + " wakafetch [options]")
	fmt.Println(ui.Clr.Bold + "Options:" + ui.Clr.Reset)

	maxWidth := 0
	for _, f := range flagDefs {
		width := len("-" + f.shortName + ", --" + f.longName + " " + f.flagType)
		if width > maxWidth {
			maxWidth = width
		}
	}

	for _, f := range flagDefs {
		flag := fmt.Sprintf("-%s, --%s", f.shortName, f.longName)
		flagLen := len(flag)
		if f.flagType != "" {
			flagLen = len(flag + " " + f.flagType)
			flag += " " + ui.Clr.Blue + f.flagType + ui.Clr.Reset
		}
		padding := strings.Repeat(" ", maxWidth-flagLen+2)
		fmt.Println("  " + ui.Clr.Green + flag + ui.Clr.Reset + padding + f.description)
	}
	os.Exit(0)
}

func loadAPIConfig(config Config) (string, string) {
	apiURL, apiKey, err := parseConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	if *config.apiKeyFlag != "" {
		apiKey = *config.apiKeyFlag
	}

	return apiURL, apiKey
}

func shouldUseSummaryAPI(config Config) bool {
	return *config.daysFlag != 0 || *config.dailyFlag
}

func handleStatsFlow(config Config, apiKey, apiURL string) {
	rangeStr := getRangeStr(*config.rangeFlag)

	data, err := fetchStats(apiKey, apiURL, rangeStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	ui.DisplayStats(data, *config.fullFlag, rangeStr)
}

func handleSummaryFlow(config Config, apiKey, apiURL string) {
	rangeStr := getRangeStr(*config.rangeFlag)
	days := *config.daysFlag

	if days == 0 {
		days = map[string]int{
			"today":         1,
			"last_7_days":   7,
			"last_30_days":  30,
			"last_6_months": 183,
			"last_year":     365,
		}[rangeStr]
	}

	data, err := fetchSummary(apiKey, apiURL, days)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *config.dailyFlag {
		if len(data.Data) == 0 {
			ui.Warnln("No daily data available")
			return
		}
		ui.DisplayBreakdown(data.Data)
		return
	}

	if len(data.Data) == 0 {
		ui.Warnln("No data available for the selected period.")
		return
	}
	ui.DisplaySummary(data, *config.fullFlag, rangeStr)
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
		log.Fatalf("Invalid range: '%s', must be one of %stoday, 7d, 30d, 6m, 1y, all", rangeFlag, ui.Clr.Green)
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
