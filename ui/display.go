package ui

import (
	"fmt"

	"github.com/sahaj-b/wakafetch/types"
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

type DisplayPayload struct {
	Heading          string
	Stats            []KV
	Languages        []types.StatItem
	Editors          []types.StatItem
	Projects         []types.StatItem
	OperatingSystems []types.StatItem
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

	statsMap := []KV{
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
		Full:             full,
	}
	render(payload)
}

func DisplaySummary(data *types.SummaryResponse, full bool, rangeStr string) {
	if data == nil || len(data.Data) == 0 {
		Warnln("No data available for the selected period: '%s'", rangeStr)
		return
	}

	languages := make(map[string]float64)
	projects := make(map[string]float64)
	editors := make(map[string]float64)
	operatingSystems := make(map[string]float64)

	for _, dayData := range data.Data {
		for _, lang := range dayData.Languages {
			languages[lang.Name] += lang.TotalSeconds
		}
		for _, proj := range dayData.Projects {
			projects[proj.Name] += proj.TotalSeconds
		}
		for _, editor := range dayData.Editors {
			editors[editor.Name] += editor.TotalSeconds
		}
		for _, os := range dayData.OperatingSystems {
			operatingSystems[os.Name] += os.TotalSeconds
		}
	}

	aggregatedLangs := mapToSortedStatItems(languages)
	aggregatedProjs := mapToSortedStatItems(projects)
	aggregatedEditors := mapToSortedStatItems(editors)
	aggregatedOS := mapToSortedStatItems(operatingSystems)

	heading := formatRangeHeading(rangeStr) + " (" + formatDateRange(data.Start, data.End) + ")"
	totalTime := timeFmt(data.CumulativeTotal.Seconds)
	dailyAvg := timeFmt(data.DailyAverage.Seconds)
	daysCoded := fmt.Sprintf("%d/%d days", data.DailyAverage.DaysMinusHolidays, data.DailyAverage.DaysIncludingHolidays)
	topProject := topItemName(aggregatedProjs, true)
	topEditor := topItemName(aggregatedEditors, false)
	topOS := topItemName(aggregatedOS, false)

	statsMap := []KV{
		{"Total Time", totalTime},
		{"Daily Avg", dailyAvg},
		{"Days Coded", daysCoded},
		{"Top Project", topProject},
		{"Top Editor", topEditor},
		{"Top OS", topOS},
	}

	payload := &DisplayPayload{
		Heading:          heading,
		Stats:            statsMap,
		Languages:        aggregatedLangs,
		Editors:          aggregatedEditors,
		Projects:         aggregatedProjs,
		OperatingSystems: aggregatedOS,
		Full:             full,
	}
	render(payload)
}
