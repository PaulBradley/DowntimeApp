package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func (app *Application) _appSetup() {
	app.exe = path.Dir(os.Args[0])
	app.ods = strings.ToLower(strings.TrimSpace(os.Getenv("ODS")))
	app.teardown = false
	app.s3timeout = 10 * time.Minute
	app.s3waiter_timeout = 2 * time.Minute
	app.LogFileOpen()
}

func (app *Application) _logAndPrint(level, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("%s %s", level, msg)

	if level == "INFO" {
		fmt.Printf("[%s] %s\n", TERMINAL_GREEN+level+TERMINAL_RESET, msg)
	} else {
		fmt.Printf("[%s] %s\n", TERMINAL_RED+level+TERMINAL_RESET, msg)
	}
}

func (app *Application) _printHeader() {
	fmt.Println()
	fmt.Println()
	fmt.Println("DowntimeApp AWS Infrastructure Provisioning")
	fmt.Println("===========================================")
	fmt.Println()
	app._logAndPrint("INFO", "Starting DowntimeApp AWS infrastructure provisioning")
}

func (app *Application) _processFlags() {
	flag.BoolVar(&app.teardown, "teardown", false, "un-provision the AWS infrastructure")
	flag.BoolVar(&app.status, "status", false, "report the AWS infrastructure set-up")
	flag.Parse()
}

func (app *Application) _startsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
