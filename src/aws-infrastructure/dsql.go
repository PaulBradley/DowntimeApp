package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dsql"
)

func (app *Application) DSQL_getDSQLClient() (context.Context, context.CancelFunc, *dsql.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), app.dsqltimeout)
	return ctx, cancel, dsql.NewFromConfig(app.GetAWSConfig(ctx))
}

func (app *Application) GetAWSConfig(ctx context.Context) (cfg aws.Config) {
	var err error

	cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(app.region))
	if err != nil {
		app._logAndPrint("ERROR", "%s %s", ERROR_AWS_CONFIG_LOAD, err.Error())
		os.Exit(1)
	}
	return cfg
}

func (app *Application) DSQL_CreateCluster(name string, deleteProtect bool) error {

	ctx, cancel, client := app.DSQL_getDSQLClient()
	defer cancel()

	input := &dsql.CreateClusterInput{
		DeletionProtectionEnabled: &deleteProtect,
		Tags: map[string]string{
			"Name": name,
		},
	}

	clusterProperties, err := client.CreateCluster(context.Background(), input)
	if err != nil {
		app._logAndPrint("ERROR", "Failed to create cluster: %v", err)
		os.Exit(1)
	}

	waiter := dsql.NewClusterActiveWaiter(client, func(o *dsql.ClusterActiveWaiterOptions) {
		o.MaxDelay = 30 * time.Second
		o.MinDelay = 10 * time.Second
		o.LogWaitAttempts = false
	})

	clusterInput := &dsql.GetClusterInput{
		Identifier: clusterProperties.Identifier,
	}

	app._logAndPrint("INFO", "Waiting for cluster %s to become ACTIVE", *clusterProperties.Arn)

	ticker := time.NewTicker(3 * time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Print(TERMINAL_GREEN + "█" + TERMINAL_RESET)
			case <-done:
				return
			}
		}
	}()

	err = waiter.Wait(ctx, clusterInput, app.dsqltimeout)
	ticker.Stop()
	close(done)
	if err != nil {
		app._logAndPrint("ERROR", "Failed waiting for cluster to become active: %v", err)
		os.Exit(1)
	}

	fmt.Println()
	app._logAndPrint("INFO", "Created Multi-AZ cluster: %s", *clusterProperties.Arn)
	return nil
}

func (app *Application) DSQL_DeleteCluster(identifier string) error {

	ctx, cancel, client := app.DSQL_getDSQLClient()
	defer cancel()

	deleteInput := &dsql.DeleteClusterInput{
		Identifier: &identifier,
	}

	result, err := client.DeleteCluster(ctx, deleteInput)
	if err != nil {
		app._logAndPrint("ERROR", "Failed to delete cluster: %v", err)
		os.Exit(1)
	}

	app._logAndPrint("INFO", "Initiated deletion of cluster: %s", *result.Arn)

	waiter := dsql.NewClusterNotExistsWaiter(client, func(options *dsql.ClusterNotExistsWaiterOptions) {
		options.MinDelay = 10 * time.Second
		options.MaxDelay = 30 * time.Second
		options.LogWaitAttempts = false
	})

	getInput := &dsql.GetClusterInput{
		Identifier: &identifier,
	}

	app._logAndPrint("INFO", "Waiting for cluster %s to be deleted...", identifier)

	ticker := time.NewTicker(3 * time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Print(TERMINAL_RED + "█" + TERMINAL_RESET)
			case <-done:
				return
			}
		}
	}()

	err = waiter.Wait(ctx, getInput, app.dsqltimeout)
	ticker.Stop()
	close(done)
	err = waiter.Wait(ctx, getInput, app.dsqltimeout)
	if err != nil {
		app._logAndPrint("ERROR", "Failed waiting for cluster to be deleted: %v", err)
		os.Exit(1)
	}

	fmt.Println()
	app._logAndPrint("INFO", "Cluster %s has been successfully deleted", identifier)
	return nil
}

func (app *Application) DSQL_GetCluster(identifier string) (clusterStatus *dsql.GetClusterOutput, err error) {

	ctx, cancel, client := app.DSQL_getDSQLClient()
	defer cancel()

	input := &dsql.GetClusterInput{
		Identifier: aws.String(identifier),
	}
	clusterStatus, err = client.GetCluster(ctx, input)

	if err != nil {
		app._logAndPrint("ERROR", "Failed to get cluster: %v", err)
		os.Exit(1)
	}

	if len(clusterStatus.Tags) == 0 {
		app._logAndPrint("INFO", "Cluster %s has no tags", identifier)
	} else {
		keys := make([]string, 0, len(clusterStatus.Tags))
		for key := range clusterStatus.Tags {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		if name, ok := clusterStatus.Tags["Name"]; ok {
			if !app.DSQL_UpdateApplicationConfig(
				name,
				aws.ToString(clusterStatus.Identifier),
				aws.ToString(clusterStatus.Endpoint)) {
				app._logAndPrint("INFO", "Skipping non DowntimeApp Cluster : %s", name)
			}
		}
	}

	return clusterStatus, nil
}

func (app *Application) DSQL_ListClusters() (clustersOutput *dsql.ListClustersOutput, err error) {

	ctx, cancel, client := app.DSQL_getDSQLClient()
	defer cancel()

	input := &dsql.ListClustersInput{}
	clustersOutput, err = client.ListClusters(ctx, input)

	if err != nil {
		app._logAndPrint("ERROR", "Failed to list clusters: %v", err)
		os.Exit(1)
	}

	for _, cluster := range clustersOutput.Clusters {
		app.DSQL_GetCluster(*cluster.Identifier)
	}

	return clustersOutput, nil
}

func (app *Application) DSQL_Provision() {
	var err error

	_, err = app.DSQL_ListClusters()
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list clusters: %v", err)
		os.Exit(1)
	}

	for index := range app.databases {
		if !app.databases[index].Found {
			app._logAndPrint("INFO", "Provisioning Cluster : %s", app.databases[index].Name)
			err = app.DSQL_CreateCluster(app.databases[index].Name, app.databases[index].DeleteProtection)
			if err != nil {
				app._logAndPrint("ERROR", "Failed to create cluster: %v", err)
				os.Exit(1)
			}
		}
		continue
	}
}

func (app *Application) DSQL_RemoveDeleteProtection(id string) (clusterStatus *dsql.UpdateClusterOutput, err error) {

	ctx, cancel, client := app.DSQL_getDSQLClient()
	defer cancel()

	input := dsql.UpdateClusterInput{
		Identifier:                &id,
		DeletionProtectionEnabled: aws.Bool(false),
	}

	clusterStatus, err = client.UpdateCluster(ctx, &input)

	if err != nil {
		app._logAndPrint("ERROR", "Failed to update cluster: %v", err)
		os.Exit(1)
	}

	app._logAndPrint("INFO", "Cluster %s updated successfully: %v", id, clusterStatus.Status)
	return clusterStatus, nil
}

func (app *Application) DSQL_Report() {
	var err error

	_, err = app.DSQL_ListClusters()
	if err != nil {
		app._logAndPrint("ERROR", "Failed to gather list of clusters: %v", err)
		os.Exit(1)
	}

	fmt.Println("-------------------------------")
	fmt.Println("C L U S T E R S   C R E A T E D")
	fmt.Println("-------------------------------")

	for index := range app.databases {
		if app.databases[index].Found {
			fmt.Printf("Cluster Name     : %s \n", app.databases[index].Name)
			fmt.Printf("Cluster Endpoint : %s\n", app.databases[index].Endpoint)
			fmt.Println("-------------------------------")
		}
		continue
	}
}

func (app *Application) DSQL_Teardown() {
	var err error

	_, err = app.DSQL_ListClusters()
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list clusters: %v", err)
		os.Exit(1)
	}

	for index := range app.databases {
		if app.databases[index].Found {
			app.DSQL_RemoveDeleteProtection(app.databases[index].Identifier)

			time.Sleep(2 * time.Second)
			app.DSQL_DeleteCluster(app.databases[index].Identifier)
		}
		continue
	}
}

func (app *Application) DSQL_UpdateApplicationConfig(name, identifier, endpoint string) bool {
	for index := range app.databases {
		if app.databases[index].Name != name {
			continue
		}

		app.databases[index].Found = true
		app.databases[index].Identifier = identifier
		app.databases[index].Endpoint = endpoint
		return true
	}

	return false
}
