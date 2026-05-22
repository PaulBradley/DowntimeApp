package main

import (
	"fmt"
	"os"
	"time"
)

type Application struct {
	exe              string
	logger           *os.File
	ods              string
	region           string
	teardown         bool
	status           bool
	s3timeout        time.Duration
	s3waiter_timeout time.Duration

	buckets   []S3_bucket
	databases []DSQL_database
}

func main() {
	app := Application{}
	app._printHeader()
	app._appSetup()
	app._processFlags()

	// define the AWS infrastructure to be
	// provisioned & monitored by the application
	app.region = "eu-west-2"
	app.databases = []DSQL_database{
		{Name: app.ods + "-downtimeapp-production"},
		{Name: app.ods + "-downtimeapp-development"},
	}
	app.setDatabaseDefaults()

	app.buckets = []S3_bucket{
		{Name: app.ods + "-downtimeapp-production"},
		{Name: app.ods + "-downtimeapp-development"},
	}
	app.setBucketDefaults()
	// E N D  O F  S E T U P

	app.ProcessFlagOverrides()
	app.Provision()
	app.Report()

	fmt.Println("== END ==")
	app.LogFileClose()
}

func (app *Application) ProcessFlagOverrides() {
	if app.status {
		app.Report()
	}

	if app.teardown {
		app.Teardown()
	}
}

func (app *Application) Provision() {
	// app.DSQL_Provision()
	app.S3_Provision()
}

func (app *Application) Report() {
	app._logAndPrint("INFO", "Gathering status details")
	time.Sleep(5 * time.Second)

	// app.DSQL_Report()
	app.S3_Report()
	app.LogFileClose()
	os.Exit(0)
}

func (app *Application) Teardown() {
	// app.DSQL_Teardown()
	app.S3_Teardown()
	app.LogFileClose()
	os.Exit(0)
}
