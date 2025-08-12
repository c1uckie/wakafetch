package ui

type Colors struct {
	Red      string
	Yellow   string
	BoldBlue string
	Blue     string
	Green    string
	Gray     string
	Bold     string
	Reset    string
}

var Clr Colors = Colors{
	Red:      "",
	Yellow:   "",
	BoldBlue: "",
	Blue:     "",
	Green:    "",
	Gray:     "",
	Bold:     "",
	Reset:    "",
}

func EnableColors() {
	Clr = Colors{
		Red:      "\033[31m",
		Yellow:   "\033[33m",
		BoldBlue: "\033[1;34m",
		Blue:     "\033[34m",
		Green:    "\033[32m",
		Gray:     "\033[90m",
		Bold:     "\033[1m",
		Reset:    "\033[0m",
	}
}
