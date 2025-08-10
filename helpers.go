package main

import (
	"fmt"
	"os"
)

func errorln(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, red+format+reset+"\n", args...)
}

func warnln(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, yellow+format+reset+"\n", args...)
}
