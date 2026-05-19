package main

func (app *Application) setDatabaseDefaults() {
	for index := range app.databases {
		app.databases[index].Found = false
		app.databases[index].DeleteProtection = true
	}
}
