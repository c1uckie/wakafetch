package ui

import (
	"fmt"

	"github.com/sahaj-b/wakafetch/types"
)

const (
	barWidth = 25
	barChar  = "🬋" // ❙ 🬋 ▆ ❘ ❚ █ ━ ▭ ╼ ━ 🬋
	spacing  = 3
)

type Field struct {
	Key string
	Val string
}

type DisplayPayload struct {
	Heading          string
	Stats            []Field
	Languages        []types.StatItem
	Editors          []types.StatItem
	Projects         []types.StatItem
	OperatingSystems []types.StatItem
	Categories       []types.StatItem
	Branches         []types.StatItem
	Machines         []types.StatItem
	DailyData        []types.DayData
	Full             bool
}

func DisplayStats(data *types.StatsResponse, full bool, rangeStr string) {
	if data == nil || (data.Data.TotalSeconds == 0 && len(data.Data.Languages) == 0 && len(data.Data.Projects) == 0) {
		Warnln("No data available for the selected period: '%s'", rangeStr)
		return
	}

	stats := data.Data
	heading := formatRangeHeading(rangeStr) + " (" + formatDateRange(stats.Start, stats.End) + ")"
	topEditor := topItemName(stats.Editors, false)
	topOS := topItemName(stats.OperatingSystems, false)
	topProject := topItemName(stats.Projects, true)

	if topProject == "unknown" {
		if len(stats.Projects) > 1 {
			topProject = stats.Projects[1].Name
		}
	}

	dailyAvg := timeFmt(stats.DailyAverage)
	totalTime := timeFmt(stats.TotalSeconds)

	statsMap := []Field{
		{"Total Time", totalTime},
		{"Daily Avg", dailyAvg},
		{"Top Project", topProject},
		{"Top Editor", topEditor},
		{"Top OS", topOS},
	}

	payload := &DisplayPayload{
		Heading:          heading,
		Stats:            statsMap,
		Languages:        stats.Languages,
		Editors:          stats.Editors,
		Projects:         stats.Projects,
		OperatingSystems: stats.OperatingSystems,
		Categories:       stats.Categories,
		Branches:         stats.Branches,
		Machines:         stats.Machines,
		DailyData:        nil,
		Full:             full,
	}
	render(payload)
}

type job struct {
	targetMap map[string]float64
	getter    func(types.DayData) []types.StatItem
}

func processJobs(data []types.DayData, jobs []job) {
	for _, dayData := range data {
		for _, j := range jobs {
			items := j.getter(dayData)
			for _, item := range items {
				j.targetMap[item.Name] += item.TotalSeconds
			}
		}
	}
}

func DisplaySummary(data *types.SummaryResponse, full bool, rangeStr string) {
	if data == nil || len(data.Data) == 0 {
		Warnln("No data available for the selected period: '%s'", rangeStr)
		return
	}

	languages := make(map[string]float64)

	// Only process additional data if full mode is on
	var projects, editors, operatingSystems, categories, branches, machines map[string]float64
	var aggregateJobs []job

	if full {
		projects = make(map[string]float64)
		editors = make(map[string]float64)
		operatingSystems = make(map[string]float64)
		categories = make(map[string]float64)
		branches = make(map[string]float64)
		machines = make(map[string]float64)

		aggregateJobs = []job{
			{languages, func(day types.DayData) []types.StatItem { return day.Languages }},
			{projects, func(day types.DayData) []types.StatItem { return day.Projects }},
			{editors, func(day types.DayData) []types.StatItem { return day.Editors }},
			{operatingSystems, func(day types.DayData) []types.StatItem { return day.OperatingSystems }},
			{categories, func(day types.DayData) []types.StatItem { return day.Categories }},
			{branches, func(day types.DayData) []types.StatItem { return day.Branches }},
			{machines, func(day types.DayData) []types.StatItem { return day.Machines }},
		}
	} else {
		// Only process languages for basic mode
		aggregateJobs = []job{
			{languages, func(day types.DayData) []types.StatItem { return day.Languages }},
		}
	}

	processJobs(data.Data, aggregateJobs)

	// Find busiest day
	busiestDay := ""
	busiestDaySeconds := 0.0
	for _, dayData := range data.Data {
		if dayData.GrandTotal.TotalSeconds > busiestDaySeconds {
			busiestDaySeconds = dayData.GrandTotal.TotalSeconds
			busiestDay = dayData.Range.Date
		}
	}

	aggregatedLangs := mapToSortedStatItems(languages)

	var aggregatedProjs, aggregatedEditors, aggregatedOS, aggregatedCategories, aggregatedBranches, aggregatedMachines []types.StatItem

	if full {
		aggregatedProjs = mapToSortedStatItems(projects)
		aggregatedEditors = mapToSortedStatItems(editors)
		aggregatedOS = mapToSortedStatItems(operatingSystems)
		aggregatedCategories = mapToSortedStatItems(categories)
		aggregatedBranches = mapToSortedStatItems(branches)
		aggregatedMachines = mapToSortedStatItems(machines)
	}

	heading := formatRangeHeading(rangeStr) + " (" + formatDateRange(data.Start, data.End) + ")"
	totalTime := timeFmt(data.CumulativeTotal.Seconds)
	dailyAvg := timeFmt(data.DailyAverage.Seconds)
	activeDays := fmt.Sprintf("%d/%d days", data.DailyAverage.DaysMinusHolidays, data.DailyAverage.DaysIncludingHolidays)

	var topProject, topEditor, topOS, topCategory string
	var numLangs, numProjects string

	if full {
		topProject = topItemName(aggregatedProjs, true)
		topEditor = topItemName(aggregatedEditors, false)
		topOS = topItemName(aggregatedOS, false)
		topCategory = topItemName(aggregatedCategories, false)
		numLangs = fmt.Sprintf("%d", len(aggregatedLangs))
		numProjects = fmt.Sprintf("%d", len(aggregatedProjs))
	}

	var statsMap []Field
	if full {
		statsMap = []Field{
			{"Total Time", totalTime},
			{"Daily Avg", dailyAvg},
			{"Active Days", activeDays},
			{"Best Day", fmt.Sprintf("%s (%s)", formatBestDay(busiestDay), timeFmt(busiestDaySeconds))},
			{"Top Project", topProject},
			{"Top Editor", topEditor},
			{"Top Category", topCategory},
			{"Top OS", topOS},
			{"Languages", numLangs},
			{"Projects", numProjects},
		}
	} else {
		statsMap = []Field{
			{"Total Time", totalTime},
			{"Daily Avg", dailyAvg},
			{"Active Days", activeDays},
			{"Best Day", fmt.Sprintf("%s (%s)", formatBestDay(busiestDay), timeFmt(busiestDaySeconds))},
		}
	}

	payload := &DisplayPayload{
		Heading:          heading,
		Stats:            statsMap,
		Languages:        aggregatedLangs,
		Editors:          aggregatedEditors,
		Projects:         aggregatedProjs,
		OperatingSystems: aggregatedOS,
		Categories:       aggregatedCategories,
		Branches:         aggregatedBranches,
		Machines:         aggregatedMachines,
		DailyData:        data.Data,
		Full:             full,
	}
	render(payload)
}
