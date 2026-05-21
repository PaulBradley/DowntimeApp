package main

import (
	"fmt"
	"log"
)

func (app *Application) _logAndPrint(level, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("%s %s", level, msg)

	if level == "INFO" {
		fmt.Printf("[%s] %s\n", TERMINAL_GREEN+level+TERMINAL_RESET, msg)
	} else {
		fmt.Printf("[%s] %s\n", TERMINAL_RED+level+TERMINAL_RESET, msg)
	}
}

func (app *Application) _startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
