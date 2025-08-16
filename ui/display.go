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
	var heading string
	if rangeStr == "all_time" {
		heading = formatRangeHeading(rangeStr)
	} else {
		heading = formatRangeHeading(rangeStr) + " (" + formatDateRange(stats.Start, stats.End) + ")"
	}

	totalTime := timeFmt(stats.TotalSeconds)
	dailyAvg := timeFmt(stats.DailyAverage)

	topProject := topItemName(stats.Projects, true)
	topEditor := topItemName(stats.Editors, false)
	topOS := topItemName(stats.OperatingSystems, false)
	topCategory := topItemName(stats.Categories, false)

	numLangs := fmt.Sprintf("%d", len(stats.Languages))
	numProjects := fmt.Sprintf("%d", len(stats.Projects))

	var statsMap []Field
	statsMap = []Field{
		{"Total Time", totalTime},
		{"Daily Avg", dailyAvg},
		{"Top Project", topProject},
		{"Top Editor", topEditor},
		{"Top Category", topCategory},
		{"Top OS", topOS},
		{"Languages", numLangs},
		{"Projects", numProjects},
	}

	payload := DisplayPayload{
		Heading:          heading,
		Stats:            statsMap,
		Languages:        stats.Languages,
		Editors:          stats.Editors,
		Projects:         stats.Projects,
		OperatingSystems: stats.OperatingSystems,
		Categories:       stats.Categories,
		Machines:         stats.Machines,
		DailyData:        nil,
		Full:             full,
	}
	render(&payload)
}

type job struct {
	targetMap map[string]float64
	getter    func(types.DayData) []types.StatItem
}

func DisplaySummary(data *types.SummaryResponse, full bool, rangeStr string) {
	if data == nil || len(data.Data) == 0 {
		Warnln("No data available for the selected period: '%s'", rangeStr)
		return
	}

	languages := make(map[string]float64)

	// Only process additional data if full mode is on
	var projects, editors, operatingSystems, categories, machines map[string]float64

	projects = make(map[string]float64)
	editors = make(map[string]float64)
	operatingSystems = make(map[string]float64)
	categories = make(map[string]float64)
	machines = make(map[string]float64)

	aggregateJobs := []job{
		{languages, func(day types.DayData) []types.StatItem { return day.Languages }},
		{projects, func(day types.DayData) []types.StatItem { return day.Projects }},
		{editors, func(day types.DayData) []types.StatItem { return day.Editors }},
		{operatingSystems, func(day types.DayData) []types.StatItem { return day.OperatingSystems }},
		{categories, func(day types.DayData) []types.StatItem { return day.Categories }},
		{machines, func(day types.DayData) []types.StatItem { return day.Machines }},
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
	aggregatedProjs := mapToSortedStatItems(projects)
	aggregatedEditors := mapToSortedStatItems(editors)
	aggregatedOS := mapToSortedStatItems(operatingSystems)
	aggregatedCategories := mapToSortedStatItems(categories)
	aggregatedMachines := mapToSortedStatItems(machines)

	heading := formatRangeHeading(rangeStr) + " (" + formatDateRange(data.Start, data.End) + ")"
	totalTime := timeFmt(data.CumulativeTotal.Seconds)
	dailyAvg := timeFmt(data.DailyAverage.Seconds)
	activeDays := fmt.Sprintf("%d/%d days", data.DailyAverage.DaysMinusHolidays, data.DailyAverage.DaysIncludingHolidays)

	topProject := topItemName(aggregatedProjs, true)
	topEditor := topItemName(aggregatedEditors, false)
	topOS := topItemName(aggregatedOS, false)
	// topCategory := topItemName(aggregatedCategories, false)

	numLangs := fmt.Sprintf("%d", len(aggregatedLangs))
	numProjects := fmt.Sprintf("%d", len(aggregatedProjs))

	var statsMap []Field
	statsMap = []Field{
		{"Total Time", totalTime},
		{"Daily Avg", dailyAvg},
		{"Active Days", activeDays},
		{"Best Day", fmt.Sprintf("%s (%s)", formatBestDay(busiestDay), timeFmt(busiestDaySeconds))},
		{"Top Project", topProject},
		{"Top Editor", topEditor},
		{"Top OS", topOS},
		{"Languages", numLangs},
		{"Projects", numProjects},
		// {"Top Category", topCategory},
	}

	payload := &DisplayPayload{
		Heading:          heading,
		Stats:            statsMap,
		Languages:        aggregatedLangs,
		Editors:          aggregatedEditors,
		Projects:         aggregatedProjs,
		OperatingSystems: aggregatedOS,
		DailyData:        data.Data,
		Full:             full,
		Categories:       aggregatedCategories,
		Machines:         aggregatedMachines,
	}
	render(payload)
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
