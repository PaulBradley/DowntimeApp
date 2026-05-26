package main

import (
	"context"
	"os"

	dta "github.com/awslabs/aurora-dsql-connectors/go/pgx/dsql"
)

func (app *Application) ApplyMigration(sql string) {

	ctx, cancel := context.WithTimeout(context.Background(), app.dsql_timeout)
	defer cancel()

	pool, err := dta.NewPool(ctx, dta.Config{
		Host: app.dsql_endpoint,
	})

	if err != nil {
		app._logAndPrint("ERROR", "Failed to create DSQL connection pool: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, sql)

	if err != nil {
		app._logAndPrint("ERROR", "Failed to execute migration %v", err)
		os.Exit(1)
	}
}
