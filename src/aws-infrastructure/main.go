package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

type DSQL_database struct {
	Name             string
	Identifier       string
	Endpoint         string
	Found            bool
	DeleteProtection bool
}

type Application struct {
	exe       string
	logger    *os.File
	teardown  bool
	region    string
	databases []DSQL_database
}

func main() {
	app := Application{}
	app.printHeader()
	app.appSetup()
	app._logAndPrint("INFO", "Starting DowntimeApp AWS infrastructure provisioning tool")

	flag.BoolVar(&app.teardown, "teardown", false, "un-provision the AWS infrastructure")
	flag.Parse()

	// define the AWS infrastructure to be
	// provisioned & monitored by the application
	app.region = "eu-west-2"
	app.databases = []DSQL_database{
		{Name: "downtimeapp-production"},
		{Name: "downtimeapp-development"},
	}
	app.setDatabaseDefaults()

	// E N D  O F  S E T U P

	if app.teardown {
		fmt.Println("== TEARDOWN ==")
		app.DSQL_Teardown()

		fmt.Println("== END ==")
		app.LogFileClose()
		os.Exit(0)
	}

	app.DSQL_Provision()

	fmt.Println("== END ==")
	app.LogFileClose()
}

func (app *Application) printHeader() {
	fmt.Println()
	fmt.Println()
	fmt.Println("DowntimeApp AWS Infrastructure Provisioning")
	fmt.Println("===========================================")
	fmt.Println()
}

func (app *Application) appSetup() {
	app.exe = path.Dir(os.Args[0])
	app.teardown = false

	app.LogFileOpen()
}
