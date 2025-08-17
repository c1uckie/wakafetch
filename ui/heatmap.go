package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/sahaj-b/wakafetch/types"
)

func heatmap(days []types.DayData) ([]string, int) {
	const heatmapChar = "■"               // █ ❐ ▪ ◼ 🙩 🙫 ⛝ ⏹ 🞕 🞔 🞖
	const highlight = "\x1b[38;2;0;%v;0m" // \x1b[38;2;R;G;Bm
	const height = 4

	if len(days) == 0 {
		return []string{}, 0
	}
	maxSecs := 0
	startDay, err := time.Parse("2006-01-02", strings.Split(days[0].Range.Start, "T")[0])
	if err != nil {
		return []string{}, 0
	}
	endDay, err := time.Parse("2006-01-02", strings.Split(days[len(days)-1].Range.Start, "T")[0])
	if err != nil {
		return []string{}, 0
	}

	for _, day := range days {
		if int(day.GrandTotal.TotalSeconds) > maxSecs {
			maxSecs = int(day.GrandTotal.TotalSeconds)
		}
	}
	output := make([]string, height)
	dataIndex := 0
	i := 0
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		strength255 := 0
		if dataIndex < len(days) && strings.Split(days[dataIndex].Range.Start, "T")[0] == d.Format("2006-01-02") {
			day := days[dataIndex]
			dataIndex++
			strength255 = int(255 * day.GrandTotal.TotalSeconds / float64(maxSecs))
		}
		char := fmt.Sprintf(highlight, strength255) + heatmapChar + "\x1b[0m"
		output[i%height] += char + " "
		i++
	}
	// trim trailing spaces
	for j := range output {
		output[j] = strings.TrimRight(output[j], " ")
	}

	// ensure same width for all lines
	for i%height != 0 {
		output[i%height] += "  "
		i++
	}
	numOfDays := int(endDay.Sub(startDay).Hours()/24) + 1
	fmt.Println("Number of days:", numOfDays)
	columns := (numOfDays + height - 1) / height // ceil
	width := columns*2 - 1
	println("Width:", width)
	return output, width
}
