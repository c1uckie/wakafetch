package ui

import "strings"

func cardify(content []string, header string, contentWidth int) ([]string, int) {
	var (
		borderTopLeft     = Clr.MidGray + "╭" + Clr.Reset
		borderTopRight    = Clr.MidGray + "╮" + Clr.Reset
		borderBottomLeft  = Clr.MidGray + "╰" + Clr.Reset
		borderBottomRight = Clr.MidGray + "╯" + Clr.Reset
		borderHorizontal  = Clr.MidGray + "─" + Clr.Reset
		borderVertical    = Clr.MidGray + "│" + Clr.Reset
	)

	if len(content) == 0 {
		return []string{}, 0
	}

	cardWidth := contentWidth + 4

	// if header is longer than content, adjust card width
	if len(header) > contentWidth {
		cardWidth = len(header) + 4
	}

	availableSpace := max(0, cardWidth-len(header)-2) // -2 for corner chars

	leftPadding := availableSpace / 2
	rightPadding := availableSpace - leftPadding

	headerLine := borderTopLeft +
		strings.Repeat(borderHorizontal, leftPadding) +
		Clr.Bold + Clr.Yellow + header + Clr.Reset +
		strings.Repeat(borderHorizontal, rightPadding) +
		borderTopRight

	result := make([]string, 0, len(content)+3)
	result = append(result, headerLine)

	// content lines
	for _, line := range content {
		padding := max(0, cardWidth-len(line)-4)
		contentLine := borderVertical + " " + line + strings.Repeat(" ", padding) + " " + borderVertical
		result = append(result, contentLine)
	}

	// bottom border
	bottomLine := borderBottomLeft + strings.Repeat(borderHorizontal, cardWidth-2) + borderBottomRight
	result = append(result, bottomLine)

	return result, cardWidth
}
