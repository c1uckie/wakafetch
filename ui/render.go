package ui

import (
	"fmt"
	"strings"
)

const (
	barWidth = 25
	barChar  = "🬋" // ❙ 🬋 ▆ ❘ ❚ █ ━ ▭ ╼ ━ 🬋
	spacing  = 1
)

type CardConfig struct {
	Title string
	Lines []string
	Width int
}

type CardSection struct {
	Left  []CardConfig
	Right []CardConfig
}

func processCardConfigs(configs []CardConfig) ([]string, int) {
	if len(configs) == 0 {
		return []string{}, 0
	}

	// Find max width across all configs
	maxWidth := 0
	for _, config := range configs {
		maxWidth = max(maxWidth, config.Width)
	}

	// Create cards using the max width
	var allCards []string
	var finalCardWidth int
	for _, config := range configs {
		rightPad := maxWidth - config.Width
		cardLines, cardWidth := cardify(config.Lines, config.Title, maxWidth, rightPad)
		finalCardWidth = cardWidth
		allCards = append(allCards, cardLines...)
	}

	return allCards, finalCardWidth
}

func renderCardSection(section CardSection) {
	leftCards, leftWidth := processCardConfigs(section.Left)
	rightCards, _ := processCardConfigs(section.Right)
	printLeftRight(leftCards, rightCards, spacing, leftWidth)
}

func render(p *DisplayPayload) {
	fields, fieldsWidth := fieldsStr(p.Heading, p.Stats)
	langLimit := len(fields)
	langGraph, langWidth := graphStr(p.Languages, langLimit)

	if p.Full {
		projectsLines, projectsWidth := graphStr(p.Projects, 0)
		categoriesLines, categoriesWidth := graphStr(p.Categories, 0)
		machinesLines, machinesWidth := graphStr(p.Machines, 0)
		editorsLines, editorsWidth := graphStr(p.Editors, 0)
		osLines, osWidth := graphStr(p.OperatingSystems, 0)

		fullSection := CardSection{
			Left: []CardConfig{
				{Title: "Languages", Lines: langGraph, Width: langWidth},
				{Title: "Projects", Lines: projectsLines, Width: projectsWidth},
				{Title: "Categories", Lines: categoriesLines, Width: categoriesWidth},
				{Title: "Machines", Lines: machinesLines, Width: machinesWidth},
			},
			Right: []CardConfig{
				{Title: "Stats", Lines: fields, Width: fieldsWidth},
				{Title: "Editors", Lines: editorsLines, Width: editorsWidth},
				{Title: "Operating Systems", Lines: osLines, Width: osWidth},
			},
		}
		renderCardSection(fullSection)

		dailyTable, tableWidth := dailyBreakdownStr(p.DailyData)
		cardTable, _ := cardify(dailyTable, "Daily Breakdown", tableWidth, 0)
		printLeftRight(cardTable, []string{}, spacing, 0)
	} else {
		langGraphCard, langWidth := cardify(langGraph, "Languages", langWidth, 0)
		printLeftRight(langGraphCard, fields, spacing, langWidth)
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

func fieldsStr(heading string, stats []Field) ([]string, int) {
	if len(stats) == 0 {
		return []string{}, 0
	}

	maxKeyLength := 0
	for _, kv := range stats {
		maxKeyLength = max(maxKeyLength, len(kv.Key))
	}

	maxWidth := len(heading)
	for _, kv := range stats {
		lineWidth := maxKeyLength + 2 + len(kv.Val)
		maxWidth = max(maxWidth, lineWidth)
	}

	output := make([]string, 0, len(stats)+2)

	headingSplit := strings.Split(heading, "(")
	headingPadded := fmt.Sprintf("%-*s", maxWidth, heading)
	headingLine := Clr.BoldBlue + headingPadded[:len(headingSplit[0])] + Clr.Reset
	if len(headingSplit) > 1 {
		remaining := headingPadded[len(headingSplit[0]):]
		headingLine += Clr.Blue + remaining + Clr.Reset
	}
	output = append(output, headingLine)

	separatorLine := fmt.Sprintf("%-*s", maxWidth, strings.Repeat("-", len(heading)))
	output = append(output, separatorLine)

	for _, kv := range stats {
		rawLine := fmt.Sprintf("%-*s%s", maxKeyLength+2, kv.Key, kv.Val)
		paddedLine := fmt.Sprintf("%-*s", maxWidth, rawLine)
		styledLine := Clr.BoldBlue + paddedLine[:maxKeyLength+2] + Clr.Reset + paddedLine[maxKeyLength+2:]
		output = append(output, styledLine)
	}

	return output, maxWidth
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
