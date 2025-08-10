package main

import (
	"fmt"
	"strings"
)

const (
	barWidth   = 25
	barChar    = "🬋" // ❙ 🬋 ▆ ❘ ❚ █ ━ ▭ ╼ ━ 🬋
	spacing    = 3
	graphLimit = 10
)

func displaySummary(data *SummaryResponse, full bool, days int) {
	if len(data.Data) == 0 {
		warnln("No data available for the selected period.")
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

	if topProject == "unknown" {
		if len(summary.Projects) > 1 {
			topProject = summary.Projects[1].Name
		}
	}

	totalTime := timeFmt(data.CumulativeTotal.Seconds, false)

	rightSide := []string{
		boldBlue + heading + reset,
		strings.Repeat("-", len(heading)),
		boldBlue + "Total Time   " + reset + totalTime,
		boldBlue + "Top Project  " + reset + topProject,
		boldBlue + "Top Editor   " + reset + topEditor,
		boldBlue + "Top OS       " + reset + topOS,
	}

	langGraph, graphWidth := getBarGraph(summary.Languages, graphLimit)

	if len(langGraph) == 0 {
		for _, line := range rightSide {
			fmt.Println(line)
		}
	} else {
		printLeftRight(langGraph, rightSide, spacing, graphWidth)
	}
	if full {
		printGraph("Editors", summary.Editors)
		printGraph("Projects", summary.Projects)
	}
}

func displayStats(data *StatsResponse, full bool, rangeStr string) {
	if data == nil || (data.Data.TotalSeconds == 0 && len(data.Data.Languages) == 0 && len(data.Data.Projects) == 0) {
		warnln("No data available for the selected period: '%s'", rangeStr)
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

	dailyAvg := timeFmt(int(stats.DailyAverage), false)

	totalTime := timeFmt(int(stats.TotalSeconds), false)

	rightSide := []string{
		boldBlue + heading + reset,
		strings.Repeat("-", len(heading)),
		boldBlue + "Total Time   " + reset + totalTime,
		boldBlue + "Daily Avg    " + reset + dailyAvg,
		boldBlue + "Top Project  " + reset + topProject,
		boldBlue + "Top Editor   " + reset + topEditor,
		boldBlue + "Top OS       " + reset + topOS,
	}

	langGraph, graphWidth := getBarGraph(stats.Languages, graphLimit)
	if len(langGraph) == 0 {
		for _, line := range rightSide {
			fmt.Println(line)
		}
	} else {
		printLeftRight(langGraph, rightSide, spacing, graphWidth)
	}

	if full {
		printGraph("Editors", stats.Editors)
		printGraph("Projects", stats.Projects)
	}
}

func formatRangeHeading(rangeStr string) string {
	lower := strings.ToLower(strings.TrimSpace(rangeStr))
	switch lower {
	case "today":
		return "Today"
	case "yesterday":
		return "Yesterday"
	case "last_7_days":
		return "Last 7 days"
	case "last_30_days":
		return "Last 30 days"
	case "last_6_months":
		return "Last 6 months"
	case "last_year":
		return "Last year"
	default:
		spaced := strings.ReplaceAll(rangeStr, "_", " ")
		if len(spaced) == 0 {
			return ""
		}
		return strings.ToUpper(spaced[:1]) + spaced[1:]
	}
}

func printGraph(title string, item []StatItem) {
	fmt.Println(bold + boldBlue + title + reset)
	graphLines, _ := getBarGraph(item, 0)
	if len(graphLines) == 0 {
		warnln("No data available for %s", title)
		return
	}
	for _, line := range graphLines {
		fmt.Println(line)
	}
}

func getBarGraph(items []StatItem, limit int) ([]string, int) {
	if len(items) == 0 {
		return []string{}, 0
	}
	count := len(items)
	if limit > 0 && limit < count {
		count = limit
	}
	visibleItems := items[:count]

	if visibleItems[0].TotalSeconds == 0 {
		return []string{}, 0
	}

	output := make([]string, 0, len(visibleItems))

	maxNameLength := 0
	for _, item := range visibleItems {
		maxNameLength = max(maxNameLength, len(item.Name))
	}
	for _, item := range visibleItems {
		if item.TotalSeconds < 60 {
			continue
		}
		barLength := int((item.TotalSeconds / visibleItems[0].TotalSeconds) * float64(barWidth))
		secondBarLength := barWidth - barLength
		if barLength < 1 {
			barLength = 1
			secondBarLength = barWidth - 1
		}
		bar := strings.Repeat(barChar, barLength)
		secondBar := strings.Repeat(barChar, secondBarLength)
		label := fmt.Sprintf("%-*s ", maxNameLength, item.Name)

		// no second bar if colors disabled
		if gray == "" {
			secondBar = strings.Repeat(" ", secondBarLength)
		}
		line := label +
			green + bar + reset +
			gray + secondBar + reset + " " +
			// green + fmt.Sprintf("%-7s", timeFmt(int(item.TotalSeconds))) + reset
			green + timeFmt(int(item.TotalSeconds), true) + reset
		output = append(output, line)
	}
	graphWidth := maxNameLength + 1 + barWidth + 1 + 7
	return output, graphWidth
}

func printLeftRight(left, right []string, spacing, leftWidth int) {
	for i, line := range left {
		if i >= len(right) {
			fmt.Println(line)
			continue
		}
		fmt.Println(line + strings.Repeat(" ", spacing) + right[i])
	}
	if len(left) < len(right) {
		pad := 0
		if len(left) > 0 {
			pad = leftWidth + spacing
		}
		for i := len(left); i < len(right); i++ {
			fmt.Println(strings.Repeat(" ", pad) + right[i])
		}
	}
}

func timeFmt(seconds int, pad bool) string {
	if seconds < 3600 {
		if pad {
			return fmt.Sprintf("%2dm %2ds", seconds/60, seconds%60)
		} else {
			return fmt.Sprintf("%dm %ds", seconds/60, seconds%60)
		}
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	if pad {
		return fmt.Sprintf("%2dh %2dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
}
