package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sahaj-b/wakafetch/types"
)

func render(p *DisplayPayload) {
	fields := fieldsStr(p.Heading, p.Stats)
	graphLimit := len(fields)
	langGraph, graphWidth := graphCard("Languages", p.Languages, graphLimit)
	printLeftRight(langGraph, fields, spacing, graphWidth)

	if p.Full {
		if len(p.Projects) > 0 || len(p.Editors) > 0 {
			fmt.Println()
			projectsGraph, width := graphCard("Projects", p.Projects, 0)
			editorsGraph, _ := graphCard("Editors", p.Editors, 0)
			printLeftRight(projectsGraph, editorsGraph, spacing, width)
		}

		if len(p.Categories) > 0 || len(p.OperatingSystems) > 0 {
			fmt.Println()
			categoriesGraph, width := graphCard("Categories", p.Categories, 0)
			osGraph, _ := graphCard("Operating Systems", p.OperatingSystems, 0)
			printLeftRight(categoriesGraph, osGraph, spacing, width)
		}

		if len(p.Machines) > 0 {
			fmt.Println()
			machinesGraph, _ := graphCard("Machines", p.Machines, 0)
			printLeftRight(machinesGraph, []string{}, spacing, 0)
		}

		if len(p.DailyData) > 0 {
			fmt.Println()
			dailyTable, width := dailyBreakdownStr(p.DailyData)
			cardTable, _ := cardify(dailyTable, "Daily Breakdown", width)
			printLeftRight(cardTable, []string{}, spacing, 0)
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

func fieldsStr(heading string, stats []Field) []string {
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
	space := strings.Repeat(" ", spacing)

	for i, line := range left {
		if i >= len(right) {
			fmt.Println(line)
			continue
		}
		fmt.Println(line + space + right[i])
	}
	if len(left) == 0 {
		spacing = 0
		leftWidth = 0
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
