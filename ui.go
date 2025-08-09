package main

import (
	"fmt"
	"strings"
)

const (
	red       = "\033[31m"
	yellow    = "\033[33m"
	blue      = "\033[1;34m"
	green     = "\033[32m"
	gray      = "\033[90m"
	bold      = "\033[1m"
	reset     = "\033[0m"
	BAR_WIDTH = 25
	BAR_CHAR  = "🬋" // ❙ 🬋 ▆ ❘ ❚ █ ━ ▭ ╼
)

func displaySummary(data *SummaryResponse, full bool, days int) {
	if len(data.Data) == 0 {
		fmt.Println(yellow + "No data available for the selected period." + reset)
		return
	}
	summary := data.Data[len(data.Data)-1]

	heading := "Today"
	if days > 1 {
		heading = fmt.Sprintf("Last %d days", days)
	}

	topEditor := "None"
	if len(summary.Editors) > 0 {
		topEditor = summary.Editors[0].Name
	}

	topOS := "None"
	if len(summary.OperatingSystems) > 0 {
		topOS = summary.OperatingSystems[0].Name
	}

	topProject := "None"
	if len(summary.Projects) > 0 {
		topProject = summary.Projects[0].Name
	}

	totalTime := timeFmt(data.CumulativeTotal.Seconds)

	rightSide := []string{
		blue + heading + reset,
		strings.Repeat("-", len(heading)),
		blue + "Total Time   " + reset + totalTime,
		blue + "Top Project  " + reset + topProject,
		blue + "Top Editor   " + reset + topEditor,
		blue + "Top OS       " + reset + topOS,
	}

	langGraph := getBarGraph(summary.Languages, 0)
	for i, line := range langGraph {
		if i >= len(rightSide) {
			fmt.Println(line)
			continue
		}
		fmt.Println(line + "   " + rightSide[i])
	}
	if len(langGraph) < len(rightSide) {
		pad := len(langGraph[0]) + 3
		for i := len(langGraph); i < len(rightSide); i++ {
			fmt.Println(strings.Repeat(" ", pad) + rightSide[i])
		}
	}
	if full {
		printGraph("Editors", getBarGraph(summary.Editors, 0))
		printGraph("Projects", getBarGraph(summary.Projects, 0))
	}
}

func displayStats(data *StatsResponse, full bool, rangeStr string) {
	if data == nil || (data.Data.TotalSeconds == 0 && len(data.Data.Languages) == 0 && len(data.Data.Projects) == 0) {
		fmt.Println(yellow + "No data available" + reset)
		return
	}

	stats := data.Data
	heading := formatRangeHeading(rangeStr)

	topEditor := "None"
	if len(stats.Editors) > 0 {
		topEditor = stats.Editors[0].Name
	}

	topOS := "None"
	if len(stats.OperatingSystems) > 0 {
		topOS = stats.OperatingSystems[0].Name
	}

	topProject := "None"
	if len(stats.Projects) > 0 {
		topProject = stats.Projects[0].Name
	}
	if topProject == "unknown" {
		if len(stats.Projects) > 1 {
			topProject = stats.Projects[1].Name
		}
	}

	totalTime := timeFmt(int(stats.TotalSeconds))

	rightSide := []string{
		blue + heading + reset,
		strings.Repeat("-", len(heading)),
		blue + "Total Time   " + reset + totalTime,
		blue + "Top Project  " + reset + topProject,
		blue + "Top Editor   " + reset + topEditor,
		blue + "Top OS       " + reset + topOS,
	}

	langGraph := getBarGraph(stats.Languages, 0)
	for i, line := range langGraph {
		if i >= len(rightSide) {
			fmt.Println(line)
			continue
		}
		fmt.Println(line + "   " + rightSide[i])
	}
	if len(langGraph) < len(rightSide) {
		pad := 0
		if len(langGraph) > 0 {
			pad = len(langGraph[0]) + 3
		}
		for i := len(langGraph); i < len(rightSide); i++ {
			fmt.Println(strings.Repeat(" ", pad) + rightSide[i])
		}
	}
	if full {
		printGraph("Editors", getBarGraph(stats.Editors, 0))
		printGraph("Projects", getBarGraph(stats.Projects, 0))
	}
}

func formatRangeHeading(rangeStr string) string {
	lower := strings.ToLower(strings.TrimSpace(rangeStr))
	switch lower {
	case "today":
		return "Today"
	case "yesterday":
		return "Yesterday"
	case "last_7_days", "last7days":
		return "Last 7 days"
	case "last_30_days":
		return "Last 30 days"
	case "last_6_months":
		return "Last 6 months"
	case "last_year":
		return "Last year"
	case "this_year":
		return "This year"
	case "week", "this_week":
		return "This week"
	case "month", "this_month":
		return "This month"
	default:
		spaced := strings.ReplaceAll(rangeStr, "_", " ")
		if len(spaced) == 0 {
			return ""
		}
		return strings.ToUpper(spaced[:1]) + spaced[1:]
	}
}

func printGraph(title string, graphLines []string) {
	fmt.Println(bold + blue + title + reset)
	for _, line := range graphLines {
		fmt.Println(line)
	}
}

func getBarGraph(items []StatItem, limit int) []string {
	output := make([]string, 0, len(items))
	if len(items) == 0 {
		return output
	}
	count := len(items)
	if limit > 0 && limit < count {
		count = limit
	}
	visibleItems := items[:count]

	if visibleItems[0].TotalSeconds == 0 {
		fmt.Println("Nothing to show")
		return output
	}

	maxNameLength := 0
	for _, item := range visibleItems {
		maxNameLength = max(maxNameLength, len(item.Name))
	}
	for _, item := range visibleItems {
		if item.TotalSeconds < 60 {
			continue
		}
		barLength := int((item.TotalSeconds / visibleItems[0].TotalSeconds) * float64(BAR_WIDTH))
		secondBarLength := BAR_WIDTH - barLength
		bar := strings.Repeat(BAR_CHAR, barLength)
		secondBar := strings.Repeat(BAR_CHAR, secondBarLength)
		label := fmt.Sprintf("%-*s ", maxNameLength, item.Name)

		line := label +
			green + bar + reset +
			gray + secondBar + reset + " " +
			green + fmt.Sprintf("%-7s", fmt.Sprintf("%dh %dm", item.Hours, item.Minutes)) + reset
		output = append(output, line)
	}
	return output
}

func timeFmt(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%dm", seconds/60)
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
