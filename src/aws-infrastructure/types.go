package main

type DSQL_database struct {
	Name             string
	DeleteProtection bool
	Endpoint         string
	Found            bool
	Identifier       string
}

type S3_bucket struct {
	ARN   string
	Found bool
	Name  string
}
