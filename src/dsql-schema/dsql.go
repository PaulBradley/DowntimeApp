package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	dta "github.com/awslabs/aurora-dsql-connectors/go/pgx/dsql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
			app._logAndPrint("ERROR", "SQL : %s", statement)
			os.Exit(1)
		}
	}
}

func (app *Application) GetTableComment(tableName string) string {
	var data strings.Builder
	var sql string = ""

	ctx, cancel := context.WithTimeout(context.Background(), app.dsql_timeout)
	defer cancel()

	sql = `SELECT COALESCE(obj_description('public.` + tableName + `'::regclass, 'pg_class'), '') AS table_comment;`

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

	for rows.Next() {
		var tableComment string
		err := rows.Scan(&tableComment)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to scan row: %v", err)
			continue
		}

		data.WriteString(tableComment + "\n\n")
	}

	if err := rows.Err(); err != nil {
		app._logAndPrint("ERROR", "Error iterating over rows: %v", err)
	}

	return data.String()
}

func (app *Application) ListTables() {

	app._logAndPrint("INFO", "Gathering table metrics")
	// time.Sleep(3 * time.Second)

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

	app.SchemaDBML()
}

func (app *Application) SchemaDBML() {
	var DBML strings.Builder
	var sql string = ""

	app._logAndPrint("INFO", "Generating DBML schema file")
	ctx, cancel := context.WithTimeout(context.Background(), app.dsql_timeout)
	defer cancel()

	sql = `
		SELECT COALESCE(table_name, '') AS table_name
		  FROM information_schema.columns
		 WHERE table_schema = 'public'
		 GROUP BY table_name
		 ORDER BY table_name`

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

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to scan row: %v", err)
			continue
		}

		DBML.WriteString("\nTable " + tableName + " {\n")

		Fields := app.GetRows(pool, `
		  SELECT
				column_name,
				data_type,
				COALESCE(character_maximum_length::text, '') AS character_maximum_length,
				COALESCE(is_nullable::text, '') AS is_nullable,
				COALESCE(is_identity::text, '') AS is_identity
				
			FROM information_schema.columns
			WHERE table_schema = 'public'
			AND table_name = '`+tableName+`'`)

		for Fields.Next() {
			var columnName string
			var dataType string
			var Length string
			var isNullable string
			var isIdentity string

			err := Fields.Scan(&columnName, &dataType, &Length, &isNullable, &isIdentity)
			if err != nil {
				app._logAndPrint("ERROR", "Failed to scan field: %v", err)
				continue
			}

			if dataType == "character" {
				dataType = "char"
			}
			if dataType == "character varying" {
				dataType = "varchar"
			}
			if strings.Contains(dataType, "timestamp") {
				dataType = "datetime"
			}
			if Length != "" {
				Length = "(" + Length + ")"
			}

			options := "["

			if len(isNullable) > 0 && isNullable == "NO" {
				options += "not null, "
			} else {
				options += "null, "
			}

			if len(isIdentity) > 0 && isIdentity == "YES" {
				options += "increment, "
			} else {
				options += "  "
			}

			Uniques := app.GetRows(pool, `
				SELECT
					COUNT(*) AS Total
				FROM
					information_schema.table_constraints tc
				JOIN
					information_schema.key_column_usage kcu
				ON
					tc.constraint_name = kcu.constraint_name
				WHERE
					tc.table_schema = 'public'
					AND tc.table_name = '`+tableName+`'
					AND kcu.column_name = '`+columnName+`'
					AND tc.constraint_type = 'PRIMARY KEY';`)

			for Uniques.Next() {
				var Total int
				err := Uniques.Scan(&Total)
				if err != nil {
					app._logAndPrint("ERROR", "Failed to scan unique constraint: %v", err)
					continue
				}

				if Total > 0 {
					options += "unique, "
				}
			}
			Uniques.Close()

			options = strings.TrimSpace(options)
			if options[len(options)-2:] == ", " {
				fmt.Println("Adding /" + options + "/")
				options += ", "
			}

			DBML.WriteString("  " + columnName + " " + dataType + " " + Length + " " + options + " note: '']\n")
		}
		Fields.Close()

		DBML.WriteString("}\n\n")
		DBML.WriteString(app.GetTableComment(tableName) + "\n\n")
	}

	if err := rows.Err(); err != nil {
		app._logAndPrint("ERROR", "Error iterating over rows: %v", err)
	}

	app.WriteDBML(DBML.String())
}

func (app *Application) WriteDBML(content string) {
	var dbml_filepath string
	dbml_filepath = app.dsql_endpoint + "-dbml.md"

	if _, err := os.Stat(dbml_filepath); err == nil {
		err = os.Remove(dbml_filepath)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to remove existing %s: %v", dbml_filepath, err)
			os.Exit(1)
		}
	}

	file, err := os.OpenFile(dbml_filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		app._logAndPrint("ERROR", "Failed to open %s: %v", dbml_filepath, err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		app._logAndPrint("ERROR", "Failed to write to %s: %v", dbml_filepath, err)
		os.Exit(1)
	}
}

func (app *Application) GetRows(pool *pgxpool.Pool, sql string) pgx.Rows {

	ctx, cancel := context.WithTimeout(context.Background(), app.dsql_timeout)
	defer cancel()

	rows, err := pool.Query(ctx, sql)
	if err != nil {
		app._logAndPrint("ERROR", "Failed to execute query: %v", err)
		os.Exit(1)
	}

	return rows
}
