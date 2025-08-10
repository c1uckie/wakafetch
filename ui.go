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

type KV struct {
	Key string
	Val string
}

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

	topEditor := topItemName(summary.Editors)

	topOS := topItemName(summary.OperatingSystems)

	topProject := topItemName(summary.Projects)

	if topProject == "unknown" {
		if len(summary.Projects) > 1 {
			topProject = summary.Projects[1].Name
		}
	}

	totalTime := timeFmt(data.CumulativeTotal.Seconds, false)

	stats := []KV{
		{"Total Time", totalTime},
		{"Top Project", topProject},
		{"Top Editor", topEditor},
		{"Top OS", topOS},
	}
	rightSide := RightSideStr(heading, stats)

	langGraph, graphWidth := graphStr(summary.Languages, graphLimit)

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

	topEditor := topItemName(stats.Editors)

	topOS := topItemName(stats.OperatingSystems)

	topProject := topItemName(stats.Projects)

	if topProject == "unknown" {
		if len(stats.Projects) > 1 {
			topProject = stats.Projects[1].Name
		}
	}

	dailyAvg := timeFmt(stats.DailyAverage, false)
	totalTime := timeFmt(stats.TotalSeconds, false)

	statsMap := []KV{
		{"Total Time", totalTime},
		{"Daily Avg", dailyAvg},
		{"Top Project", topProject},
		{"Top Editor", topEditor},
		{"Top OS", topOS},
	}
	rightSide := RightSideStr(heading, statsMap)

	langGraph, graphWidth := graphStr(stats.Languages, graphLimit)
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
	fmtRangeMap := map[string]string{
		"today":         "Today",
		"yesterday":     "Yesterday",
		"last_7_days":   "Last 7 days",
		"last_30_days":  "Last 30 days",
		"last_6_months": "Last 6 months",
		"last_year":     "Last year",
		"all_time":      "All time",
	}
	if val, exists := fmtRangeMap[lower]; exists {
		return val
	}
	return strings.ToUpper(lower[:1]) + lower[1:]
}

func printGraph(title string, item []StatItem) {
	fmt.Println(bold + boldBlue + title + reset)
	graphLines, _ := graphStr(item, 0)
	if len(graphLines) == 0 {
		warnln("No data available for %s", title)
		return
	}
	for _, line := range graphLines {
		fmt.Println(line)
	}
}

func graphStr(items []StatItem, limit int) ([]string, int) {
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
			green + timeFmt(item.TotalSeconds, true) + reset
		output = append(output, line)
	}
	graphWidth := maxNameLength + 1 + barWidth + 1 + 7
	return output, graphWidth
}

func RightSideStr(heading string, stats []KV) []string {
	if len(stats) == 0 {
		return []string{}
	}

	maxKeyLength := 0
	for _, kv := range stats {
		maxKeyLength = max(maxKeyLength, len(kv.Key))
	}

	output := make([]string, 0, len(stats)+2)
	output = append(output, boldBlue+heading+reset)
	output = append(output, strings.Repeat("-", len(heading)))
	for _, kv := range stats {
		line := boldBlue + fmt.Sprintf("%-*s", maxKeyLength+2, kv.Key) + reset + kv.Val
		output = append(output, line)
	}
	return output
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

func timeFmt(seconds float64, pad bool) string {
	sec := int(seconds)
	if sec < 3600 {
		if pad {
			return fmt.Sprintf("%2dm %2ds", sec/60, sec%60)
		} else {
			return fmt.Sprintf("%dm %ds", sec/60, sec%60)
		}
	}
	hours := sec / 3600
	minutes := (sec % 3600) / 60
	if pad {
		return fmt.Sprintf("%2dh %2dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
}

func topItemName(items []StatItem) string {
	if len(items) == 0 {
		return "None"
	}
	return items[0].Name
}
