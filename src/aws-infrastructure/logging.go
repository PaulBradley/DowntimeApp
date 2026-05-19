package main

import (
	"log"
	"os"
)

func (app *Application) LogFileOpen() {
	var err error
	app.logger, err = os.OpenFile(app.exe+"/downtimeapp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		println("ERROR:" + err.Error())
		os.Exit(1)
	}
	log.SetOutput(app.logger)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (app *Application) LogFileClose() {
	app.logger.Sync()
	app.logger.Close()
}
