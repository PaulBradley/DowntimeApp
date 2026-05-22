package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (app *Application) S3_getS3Client() (context.Context, context.CancelFunc, *s3.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), app.s3timeout)
	return ctx, cancel, s3.NewFromConfig(app.GetAWSConfig(ctx))
}

func (app *Application) S3_CreateBucket(name string) error {

	ctx, cancel, s3Client := app.S3_getS3Client()
	defer cancel()

	input := &s3.CreateBucketInput{
		Bucket: &name,
	}

	if app.region != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(app.region),
		}
	}

	// create bucket
	_, err := s3Client.CreateBucket(ctx, input)
	if err != nil {
		return fmt.Errorf("Failed to create bucket. %v", err)
	}

	waiter := s3.NewBucketExistsWaiter(s3Client)
	if err := waiter.Wait(ctx, &s3.HeadBucketInput{Bucket: &name}, app.s3waiter_timeout); err != nil {
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

func (app *Application) S3_DeleteBucket(name string) error {

	ctx, cancel, s3Client := app.S3_getS3Client()
	defer cancel()

	_, err := s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &name,
	})
	if err != nil {
		app._logAndPrint("ERROR", "Failed to delete bucket: %v", err)
		os.Exit(1)
	}

	waiter := s3.NewBucketNotExistsWaiter(s3Client)
	if err := waiter.Wait(ctx, &s3.HeadBucketInput{Bucket: &name}, app.s3waiter_timeout); err != nil {
		app._logAndPrint("ERROR", "Failed waiting for bucket to be deleted: %v", err)
		os.Exit(1)
	}

	app._logAndPrint("INFO", "Deleted Bucket: %s", name)
	return nil
}

func (app *Application) S3_ListBuckets() (buckets []string, err error) {

	ctx, cancel, s3Client := app.S3_getS3Client()
	defer cancel()

	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list buckets: %v", err)
		os.Exit(1)
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

	_, err = app.S3_ListBuckets()
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list buckets: %v", err)
		os.Exit(1)
	}

	for index := range app.buckets {
		if !app.buckets[index].Found {
			app._logAndPrint("INFO", "Provisioning Bucket: %s", app.buckets[index].Name)
			err = app.S3_CreateBucket(app.buckets[index].Name)
			if err != nil {
				app._logAndPrint("ERROR", "Failed to create Bucket: %v", err)
				os.Exit(1)
			}
		}
		continue
	}
}

func (app *Application) S3_PurgeObjects(name string) {

	ctx, cancel, s3Client := app.S3_getS3Client()
	defer cancel()

	deleteObject := func(bucket, key, versionId *string) {
		ctx, cancel, s3Client := app.S3_getS3Client()
		defer cancel()

		app._logAndPrint("INFO", "Deleting object: %s/%s", *key, aws.ToString(versionId))
		_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket:    bucket,
			Key:       key,
			VersionId: versionId,
		})
		if err != nil {
			app._logAndPrint("ERROR", "Failed to delete object: %v", err)
			os.Exit(1)
		}
	}

	in := &s3.ListObjectsV2Input{Bucket: &name}
	for {
		out, err := s3Client.ListObjectsV2(ctx, in)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to list objects: %v", err)
			os.Exit(1)
		}

		for _, item := range out.Contents {
			deleteObject(&name, item.Key, nil)
		}

		if *out.IsTruncated {
			in.ContinuationToken = out.ContinuationToken
		} else {
			break
		}
	}

	inVer := &s3.ListObjectVersionsInput{Bucket: &name}
	for {
		out, err := s3Client.ListObjectVersions(ctx, inVer)
		if err != nil {
			app._logAndPrint("ERROR", "Failed to list version objects: %v", err)
			os.Exit(1)
		}

		for _, item := range out.DeleteMarkers {
			deleteObject(&name, item.Key, item.VersionId)
		}

		for _, item := range out.Versions {
			deleteObject(&name, item.Key, item.VersionId)
		}

		if *out.IsTruncated {
			inVer.VersionIdMarker = out.NextVersionIdMarker
			inVer.KeyMarker = out.NextKeyMarker
		} else {
			break
		}
	}
}

func (app *Application) S3_Report() {
	var err error

	_, err = app.S3_ListBuckets()
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list buckets: %v", err)
		os.Exit(1)
	}

	fmt.Println("-------------------------------")
	fmt.Println("B U C K E T S   C R E A T E D  ")
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

func (app *Application) S3_Teardown() {
	var err error

	_, err = app.S3_ListBuckets()
	if err != nil {
		app._logAndPrint("ERROR", "Failed to list buckets: %v", err)
		os.Exit(1)
	}

	for index := range app.buckets {
		if app.buckets[index].Found {
			app._logAndPrint("INFO", "Purging objects from Bucket : %s", app.buckets[index].Name)
			app.S3_PurgeObjects(app.buckets[index].Name)
			app.S3_DeleteBucket(app.buckets[index].Name)
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
