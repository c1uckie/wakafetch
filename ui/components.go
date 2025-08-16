package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sahaj-b/wakafetch/types"
)

func render(p *DisplayPayload) {
	rightSide := rightSideStr(p.Heading, p.Stats)
	graphLimit := len(rightSide)
	langGraph, graphWidth := graphStr(p.Languages, graphLimit)
	printLeftRight(langGraph, rightSide, spacing, graphWidth)

	if p.Full {
		divider := "\n" + strings.Repeat("─", 100) + "\n"

		if len(p.Projects) > 0 || len(p.Editors) > 0 {
			fmt.Println(Clr.Gray + divider + Clr.Reset)
			projectsGraph, maxWidth := graphWithHeaderStr("Projects", p.Projects)
			editorsGraph, _ := graphWithHeaderStr("Editors", p.Editors)
			printLeftRight(projectsGraph, editorsGraph, spacing, maxWidth)
		}

		if len(p.Categories) > 0 || len(p.OperatingSystems) > 0 {
			fmt.Println(Clr.Gray + divider + Clr.Reset)
			categoriesGraph, maxWidth := graphWithHeaderStr("Categories", p.Categories)
			osGraph, _ := graphWithHeaderStr("Operating Systems", p.OperatingSystems)
			printLeftRight(categoriesGraph, osGraph, spacing, maxWidth)
		}

		if len(p.Machines) > 0 {
			fmt.Println(Clr.Gray + divider + Clr.Reset)
			machinesGraph, _ := graphWithHeaderStr("Machines", p.Machines)
			printLeftRight(machinesGraph, []string{}, spacing, 0)
		}

		if len(p.DailyData) > 0 {
			fmt.Println(Clr.Gray + divider + Clr.Reset)
			dailyTable := dailyBreakdownStr(p.DailyData)
			printLeftRight(dailyTable, []string{}, spacing, 0)
		}
	}
}

func mapToSortedStatItems(m map[string]float64) []types.StatItem {
	items := make([]types.StatItem, 0, len(m))
	for name, seconds := range m {
		items = append(items, types.StatItem{Name: name, TotalSeconds: seconds})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].TotalSeconds > items[j].TotalSeconds
	})
	return items
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

func graphWithHeaderStr(title string, item []types.StatItem) ([]string, int) {
	graphLines, width := graphStr(item, 0)

	if len(graphLines) == 0 {
		return []string{}, 0
	}

	output := make([]string, 0, len(item)+2)
	output = append(output, Clr.Yellow+Clr.Bold+fmt.Sprintf("%-*s", width, title)+Clr.Reset)
	output = append(output, graphLines...)
	return output, width
}

func graphStr(items []types.StatItem, limit int) ([]string, int) {
	if len(items) == 0 {
		return []string{}, 0
	}
	count := len(items)
	if limit > 0 && limit < count {
		count = limit
	}
	visibleItems := items[:count]

	if len(visibleItems) == 0 || visibleItems[0].TotalSeconds == 0 {
		return []string{}, 0
	}

	output := make([]string, 0, len(visibleItems))

	maxNameLength := 0
	maxSeconds := 0.0
	for _, item := range visibleItems {
		maxNameLength = max(maxNameLength, len(item.Name))
		maxSeconds = max(maxSeconds, item.TotalSeconds)
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

		if Clr.Gray == "" {
			secondBar = strings.Repeat(" ", secondBarLength)
		}
		line := label +
			Clr.Green + bar + Clr.Reset +
			Clr.Gray + secondBar + Clr.Reset + " " +
			Clr.Green + timeFmtPad(item.TotalSeconds, maxSeconds) + Clr.Reset
		output = append(output, line)
	}
	graphWidth := maxNameLength + 1 + barWidth + 1 + len(timeFmtPad(maxSeconds, maxSeconds))
	return output, graphWidth
}

func rightSideStr(heading string, stats []Field) []string {
	if len(stats) == 0 {
		return []string{}
	}

	maxKeyLength := 0
	for _, kv := range stats {
		maxKeyLength = max(maxKeyLength, len(kv.Key))
	}

	output := make([]string, 0, len(stats)+2)
	headingSplit := strings.Split(heading, "(")
	headingStr := Clr.BoldBlue + heading + Clr.Reset
	if len(headingSplit) > 1 {
		headingStr = Clr.BoldBlue + headingSplit[0] + Clr.Reset + Clr.Blue + "(" + headingSplit[1] + Clr.Reset
	}
	output = append(output, headingStr)
	output = append(output, strings.Repeat("-", len(heading)))
	for _, kv := range stats {
		line := Clr.BoldBlue + fmt.Sprintf("%-*s", maxKeyLength+2, kv.Key) + Clr.Reset + kv.Val
		output = append(output, line)
	}
	return output
}

func printLeftRight(left, right []string, spacing, leftWidth int) {
	separator := ""
	space := strings.Repeat(" ", spacing)
	if len(right) > 0 {
		separator = space + Clr.Gray + "│" + Clr.Reset + space
	}

	for i, line := range left {
		if i >= len(right) {
			fmt.Println(line + separator)
			continue
		}
		fmt.Println(line + separator + right[i])
	}
	if len(left) == 0 {
		spacing = 0
		leftWidth = 0
		separator = ""
	}
	if len(left) < len(right) {
		pad := 0
		if len(left) > 0 {
			pad = leftWidth + spacing
		}
		for i := len(left); i < len(right); i++ {
			fmt.Println(strings.Repeat(" ", pad) + separator + right[i])
		}
	}
}
