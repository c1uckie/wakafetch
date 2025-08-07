package main

import (
	"fmt"
)

func prettyPrint(data *SummaryResponse, full bool) {
	if full {
		fmt.Println("\nFull Statistics:")
	}
}
