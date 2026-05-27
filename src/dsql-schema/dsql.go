package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	dta "github.com/awslabs/aurora-dsql-connectors/go/pgx/dsql"
	"github.com/olekukonko/tablewriter"
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

	statements := strings.Split(sql, ";")
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err = pool.Exec(ctx, statement)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to execute migration statement %d: %v", i+1, err)
			os.Exit(1)
		}
	}
}

func (app *Application) ListTables() {

	var sql = `
		SELECT
			t.tablename,
			c.reltuples::bigint AS estimated_rows
		FROM pg_catalog.pg_tables t
		JOIN pg_catalog.pg_class c
		  ON c.relname = t.tablename
		 AND c.relnamespace =
		 	(
				SELECT oid FROM pg_catalog.pg_namespace n WHERE n.nspname = t.schemaname
			)
	   WHERE t.schemaname = 'public'
		 AND t.schemaname != 'information_schema'
	   ORDER BY estimated_rows DESC;`

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

	rows, err := pool.Query(ctx, sql)

	if err != nil {
		app._logAndPrint("ERROR", "Failed to query tables %v", err)
		os.Exit(1)
	}
	defer rows.Close()

	data := [][]string{
		{"TABLE NAME", "ROW COUNT (ESTIMATE)"},
	}

	for rows.Next() {
		var tableName string
		var rowCount int64

		err := rows.Scan(&tableName, &rowCount)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to scan row: %v", err)
			continue
		}

		data = append(data, []string{tableName, strconv.FormatInt(rowCount, 10)})
	}

	if err := rows.Err(); err != nil {
		app._logAndPrint("ERROR", "Error iterating over rows: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	table.Bulk(data[1:])
	table.Render()
}
