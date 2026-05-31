package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (app *Application) _appSetup() {
	app.exe = path.Dir(os.Args[0])
	app.ods = strings.ToLower(strings.TrimSpace(os.Getenv("ODS")))
	app.dsql_timeout = 10 * time.Minute

}

func (app *Application) _generateEnvironments() {
	cellar := "01KSM3WRPK2D9K04RS92MBYSHT"
	envs := []string{"PRODUCTION", "MOCK", "CERT", "TRAINING", "BUILD"}
	salt := os.Getenv("SALT")

	for _, env := range envs {
		vault := fmt.Sprintf("%s%s%s", cellar, env, salt)
		hash := strings.ToUpper(fmt.Sprintf("%x", sha512.Sum512([]byte(vault))))
		fmt.Printf("INSERT INTO ENVIRONMENTS (cellar, vault, environment, is_enabled) VALUES ('%s', '%s', '%s', 'Y');\n", cellar, hash[:26], env)
	}
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
	fmt.Println("DowntimeApp AWS Schema Migration Utility")
	fmt.Println("========================================")
	fmt.Println()
	app._logAndPrint("INFO", "Starting the DowntimeApp schema migration utility")
	app._logAndPrint("INFO", "Version (Git Commit Hash): %s", GIT_COMMIT_HASH)
}

func (app *Application) _printDSQLEndpoint() {
	if app.production {
		app.dsql_endpoint = strings.TrimSpace(os.Getenv(app.ods + "-downtimeapp-production"))
	} else {
		app.dsql_endpoint = strings.TrimSpace(os.Getenv(app.ods + "-downtimeapp-development"))
	}
	app._logAndPrint("INFO", "Using DSQL Endpoint: %s", app.dsql_endpoint)
}

func (app *Application) _printMigrationMethod() {
	if app.rollback {
		app._logAndPrint("INFO", "Rolling back to migration ID: %d", app.migration_id)
	} else {
		app._logAndPrint("INFO", "Upgrading to migration ID: %d", app.migration_id)
	}
}

func (app *Application) _processFlags() {
	flag.BoolVar(&app.environments, "environments", false, "generate vault ids for a given cellar")
	flag.BoolVar(&app.list_tables, "list-tables", false, "list tables and their row counts")
	flag.BoolVar(&app.production, "production", false, "run the migration in production mode")
	flag.BoolVar(&app.rollback, "rollback", false, "rollback the schema migrations")
	flag.IntVar(&app.migration_id, "id", 0, "the migration id to upgrade/rollback to/from")
	flag.Parse()
}

var fileRe = regexp.MustCompile(`^(\d+)(?:_.*)?\.sql$`)

func (app *Application) readMigrationFiles(dir string, max int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		app._logAndPrint("ERROR", "Failed %s", err.Error())
		os.Exit(1)
	}

	type entry struct {
		num  int
		path string
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		m := fileRe.FindStringSubmatch(name)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n <= max {
			app.migrations = append(app.migrations, Migration{id: n, filename: filepath.Join(dir, name)})
		}
	}
}
