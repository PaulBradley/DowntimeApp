package main

import (
	"os"
	"sort"
	"strings"
	"time"
)

type Migration struct {
	id       int
	filename string
}

type Application struct {
	dsql_endpoint string
	dsql_timeout  time.Duration
	exe           string
	logger        *os.File
	list_tables   bool
	migration_id  int
	ods           string
	production    bool
	region        string
	rollback      bool
	migrations    []Migration
}

var GIT_COMMIT_HASH string

func main() {
	app := Application{}
	app._appSetup()
	app._logFileOpen()
	app._processFlags()
	app._printHeader()
	app._printDSQLEndpoint()

	if app.list_tables {
		app.ListTables()
		os.Exit(0)
	}

	app._printMigrationMethod()
	app.readMigrationFiles("./migrations/up/", 999)

	if app.rollback {
		sort.Slice(app.migrations, func(i, j int) bool {
			return app.migrations[i].id > app.migrations[j].id
		})
	}

	for _, migration := range app.migrations {
		if app.rollback && migration.id < app.migration_id {
			continue
		}
		if !app.rollback && migration.id > app.migration_id {
			continue
		}

		script_filename := ""
		if app.rollback {
			script_filename = strings.Replace(migration.filename, "migrations/up/", "migrations/down/", 1)
		} else {
			script_filename = migration.filename
		}

		migration_script, err := os.ReadFile(script_filename)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to read migration file %s: %v", script_filename, err)
			os.Exit(1)
		}

		app.ApplyMigration(string(migration_script))
		app._logAndPrint("INFO", "Applied migration: %s", script_filename)
	}

	app._logFileClose()
}
