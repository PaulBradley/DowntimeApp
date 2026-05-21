package main

func (app *Application) setDatabaseDefaults() {
	for index := range app.databases {
		app.databases[index].Found = false
		app.databases[index].DeleteProtection = true
	}
}

func (app *Application) setBucketDefaults() {
	for index := range app.buckets {
		app.buckets[index].Found = false
	}
}
