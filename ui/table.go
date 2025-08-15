package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sahaj-b/wakafetch/types"
)

func printHeader(title string) {
	fmt.Println(Clr.Bold + Clr.Yellow + title + Clr.Reset)
}

func printDailyBreakdown(dailyData []types.DayData) {
	if len(dailyData) == 0 {
		return
	}

	// sort descending by date
	sortedDays := make([]types.DayData, len(dailyData))
	copy(sortedDays, dailyData)

	sort.Slice(sortedDays, func(i, j int) bool {
		dateI := strings.Split(sortedDays[i].Range.Start, "T")[0]
		dateJ := strings.Split(sortedDays[j].Range.Start, "T")[0]
		return dateI > dateJ
	})

	printHeader("Daily Breakdown")

	maxSecs := findMaxDailySeconds(sortedDays)
	cols := calculateDailyColumnWidths(sortedDays)

	printDailyHeaders(cols)
	printDailySeparator(cols)
	printDailyRows(sortedDays, cols, maxSecs)
}

func findMaxDailySeconds(dailyData []types.DayData) float64 {
	maxSecs := 0.0
	for _, day := range dailyData {
		if day.GrandTotal.TotalSeconds > maxSecs {
			maxSecs = day.GrandTotal.TotalSeconds
		}
	}
	return maxSecs
}

type dailyColumns struct {
	date    int
	total   int
	lang    int
	project int
}

func calculateDailyColumnWidths(dailyData []types.DayData) dailyColumns {
	cols := dailyColumns{
		date:    4, // "Date"
		total:   5, // "Total"
		lang:    8, // "Language"
		project: 7, // "Project"
	}

	maxSecs := findMaxDailySeconds(dailyData)

	for _, day := range dailyData {
		if day.GrandTotal.TotalSeconds < 60 {
			continue
		}

		dateStr := formatDailyDate(day.Range.Start)
		totalStr := timeFmtPad(day.GrandTotal.TotalSeconds, maxSecs)
		topLang := topItemName(day.Languages, false)
		topProj := topItemName(day.Projects, true)

		if len(dateStr) > cols.date {
			cols.date = len(dateStr)
		}
		if len(totalStr) > cols.total {
			cols.total = len(totalStr)
		}
		if len(topLang) > cols.lang {
			cols.lang = len(topLang)
		}
		if len(topProj) > cols.project {
			cols.project = len(topProj)
		}
	}

	// Add space for bars (max 10 chars + 1 space)
	cols.total += 11

	return cols
}

func printDailyHeaders(cols dailyColumns) {
	headerDate := fmt.Sprintf("%-*s", cols.date, "Date")
	headerTotal := fmt.Sprintf("%-*s", cols.total, "Total")
	headerLang := fmt.Sprintf("%-*s", cols.lang, "Language")
	headerProj := fmt.Sprintf("%-*s", cols.project, "Project")

	fmt.Println(Clr.Blue+headerDate+Clr.Reset, "│", Clr.Blue+headerTotal+Clr.Reset, "│", Clr.Blue+headerLang+Clr.Reset, "│", Clr.Blue+headerProj+Clr.Reset)
}

func printDailySeparator(cols dailyColumns) {
	separator := strings.Repeat("─", cols.date) + "─┼─" +
		strings.Repeat("─", cols.total) + "─┼─" +
		strings.Repeat("─", cols.lang) + "─┼─" +
		strings.Repeat("─", cols.project)
	fmt.Println(separator)
}

func printDailyRows(dailyData []types.DayData, cols dailyColumns, maxSecs float64) {
	for _, day := range dailyData {
		if day.GrandTotal.TotalSeconds < 60 {
			continue
		}

		date := fmt.Sprintf("%-*s", cols.date, formatDailyDate(day.Range.Start))

		barLength := int((day.GrandTotal.TotalSeconds / maxSecs) * 10)
		if barLength < 1 && day.GrandTotal.TotalSeconds > 0 {
			barLength = 1
		}
		bar := strings.Repeat("🬋", barLength)
		timeStr := timeFmtPad(day.GrandTotal.TotalSeconds, maxSecs)
		timeWithBar := timeStr + " " + bar
		totalFormatted := fmt.Sprintf("%-*s", cols.total, timeWithBar)

		topLang := fmt.Sprintf("%-*s", cols.lang, topItemName(day.Languages, false))
		topProj := fmt.Sprintf("%-*s", cols.project, topItemName(day.Projects, true))

		fmt.Println(date, "│", Clr.Green+totalFormatted+Clr.Reset, "│", topLang, "│", topProj)
	}
}
