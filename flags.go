package main

import (
	"flag"
	"fmt"
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
	heatmapFlag *bool
	noColorFlag *bool
	helpFlag    *bool
}

type flagInfo struct {
	longName    string
	shortName   string
	defaultVal  any
	description string
	flagType    string
}

func parseFlags() Config {
	config := Config{}
	registeredFlags = nil

	config.rangeFlag = config.stringFlag("range", "r", "7d", "Range of data to fetch (today/7d/30d/6m/1y/all) (default: 7d)")
	config.daysFlag = config.intFlag("days", "d", 0, "Number of days to fetch data for (overrides --range)")
	config.fullFlag = config.boolFlag("full", "f", false, "Display full statistics")
	config.dailyFlag = config.boolFlag("daily", "D", false, "Display daily breakdown")
	config.heatmapFlag = config.boolFlag("heatmap", "H", false, "Display heatmap of daily activity")
	config.apiKeyFlag = config.stringFlag("api-key", "k", "", "Your WakaTime/Wakapi API key (overrides config)")
	config.noColorFlag = config.boolFlag("no-colors", "n", false, "Disable colored output")
	config.helpFlag = config.boolFlag("help", "h", false, "Display help information")

	flag.Usage = showCustomHelp
	flag.Parse()

	if *config.noColorFlag || !colorsShouldBeEnabled() {
		ui.DisableColors()
	}

	if *config.helpFlag {
		showCustomHelp()
	}

	if *config.daysFlag < 0 {
		ui.Errorln("Invalid value for --days: must be a positive integer")
	}

	return config
}

var registeredFlags []flagInfo

func (c *Config) stringFlag(long, short, def, desc string) *string {
	registeredFlags = append(registeredFlags, flagInfo{long, short, def, desc, "string"})
	val := flag.String(long, def, "")
	flag.StringVar(val, short, def, "")
	return val
}

func (c *Config) boolFlag(long, short string, def bool, desc string) *bool {
	registeredFlags = append(registeredFlags, flagInfo{long, short, def, desc, ""})
	val := flag.Bool(long, def, "")
	flag.BoolVar(val, short, def, "")
	return val
}

func (c *Config) intFlag(long, short string, def int, desc string) *int {
	registeredFlags = append(registeredFlags, flagInfo{long, short, def, desc, "int"})
	val := flag.Int(long, def, "")
	flag.IntVar(val, short, def, "")
	return val
}

func showCustomHelp() {
	fmt.Println(ui.Clr.Bold + "Usage:" + ui.Clr.Reset + " wakafetch [options]")
	fmt.Println(ui.Clr.Bold + "Options:" + ui.Clr.Reset)

	maxWidth := 0
	for _, f := range registeredFlags {
		width := len("-" + f.shortName + ", --" + f.longName + " " + f.flagType)
		if width > maxWidth {
			maxWidth = width
		}
	}

	for _, f := range registeredFlags {
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
