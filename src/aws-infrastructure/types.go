package main

type DSQL_database struct {
	Name             string
	DeleteProtection bool
	Endpoint         string
	Found            bool
	Identifier       string
}

type S3_bucket struct {
	ARN               string
	DeleteProtection  bool
	Encryption        string
	Found             bool
	LoggingEnabled    bool
	Name              string
	VersioningEnabled bool
}
