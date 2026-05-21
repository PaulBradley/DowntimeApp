package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (app *Application) S3_CreateBucket(ctx context.Context, region string, name string) error {

	cfg := app.GetAWSConfig(ctx, region)

	s3Client := s3.NewFromConfig(cfg)
	input := &s3.CreateBucketInput{
		Bucket: &name,
	}

	if region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		}
	}

	// create bucket
	_, err := s3Client.CreateBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("Failed to create bucket. %v", err)
	}

	waiter := s3.NewBucketExistsWaiter(s3Client)
	if err := waiter.Wait(ctx, &s3.HeadBucketInput{Bucket: &name}, 2*time.Minute); err != nil {
		return fmt.Errorf("Failed waiting for bucket to exist: %w", err)
	}

	// disable versioning
	_, err = s3Client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
		Bucket: &name,
		VersioningConfiguration: &types.VersioningConfiguration{
			Status: types.BucketVersioningStatusSuspended,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to suspend bucket versioning: %w", err)
	}

	// enable encryption
	_, err = s3Client.PutBucketEncryption(ctx, &s3.PutBucketEncryptionInput{
		Bucket: &name,
		ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
			Rules: []types.ServerSideEncryptionRule{
				{
					ApplyServerSideEncryptionByDefault: &types.ServerSideEncryptionByDefault{
						SSEAlgorithm: types.ServerSideEncryptionAes256,
					},
					BucketKeyEnabled: aws.Bool(false),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to configure bucket encryption: %w", err)
	}

	return nil
}

func (app *Application) S3_ListBuckets(ctx context.Context, region string) (buckets []string, err error) {

	cfg := app.GetAWSConfig(ctx, region)

	s3Client := s3.NewFromConfig(cfg)

	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("Failed to list buckets: %w", err)
	}

	buckets = make([]string, 0, len(resp.Buckets))
	for _, b := range resp.Buckets {
		if b.Name != nil && app.ods != "" && !app._startsWith(*b.Name, app.ods) {
			continue
		}
		if b.Name != nil {
			buckets = append(buckets, *b.Name)
			if b.BucketArn != nil {
				app.S3_Update_Bucket(*b.Name, *b.BucketArn)
			}
		}
	}
	return buckets, nil
}

func (app *Application) S3_Provision() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	_, err = app.S3_ListBuckets(ctx, app.region)
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list buckets: %v", err)
		os.Exit(1)
	}

	for index := range app.buckets {
		if !app.buckets[index].Found {
			app._logAndPrint("INFO", "Provisioning Bucket : %s", app.buckets[index].Name)
			err = app.S3_CreateBucket(ctx, app.region, app.buckets[index].Name)
			if err != nil {
				app._logAndPrint("ERROR", "Failed to create Bucket: %v", err)
				os.Exit(1)
			}
		}
		continue
	}

	time.Sleep(5 * time.Second)
	app.S3_Report()
}

func (app *Application) S3_Report() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	_, err = app.S3_ListBuckets(ctx, app.region)
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list buckets: %v", err)
		os.Exit(1)
	}

	fmt.Println("-------------------------------")
	fmt.Println("B U C K E T S   C R E A T E D")
	fmt.Println("-------------------------------")

	for index := range app.buckets {
		if app.buckets[index].Found {
			fmt.Printf("Bucket Name : %s \n", app.buckets[index].Name)
			fmt.Printf("Bucket ARN  : %s\n", app.buckets[index].ARN)
			fmt.Println("-------------------------------")
		}
		continue
	}
}

func (app *Application) S3_Update_Bucket(name, arn string) bool {
	for index := range app.buckets {
		if app.buckets[index].Name != name {
			continue
		}

		app.buckets[index].Found = true
		app.buckets[index].ARN = arn
		return true
	}

	return false
}
