package s3

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// Function signature for those performing the checks
type CheckFn func() error

var (
	client *s3.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client = s3.NewFromConfig(cfg)
}

// Check checks for best practices for AWS S3
func Check() error {
	fmt.Println("Checking S3...")
	checks := []CheckFn{
		CheckPublicAccessBlockingEnabled,
	}

	for _, check := range checks {
		err := check()
		if err != nil {
			return err
		}
	}
	return nil
}

// getAllS3Buckets returns a list of all S3 buckets
func getAllS3Buckets() ([]types.Bucket, error) {
	out, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	return out.Buckets, nil
}

// CheckBlockingPublicAccess checks if public access is enabled, which may expose
// the bucket to the public.
func CheckPublicAccessBlockingEnabled() error {
	buckets, err := getAllS3Buckets()
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		out, err := client.GetPublicAccessBlock(context.TODO(), &s3.GetPublicAccessBlockInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			// attempt to interpret AWS-specific errors
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				code := apiErr.ErrorCode()
				message := apiErr.ErrorMessage()

				if code == "NoSuchPublicAccessBlockConfiguration" {
					fmt.Println(*bucket.Name, "- does not have public access block configuration")
					continue
				}

				fmt.Printf("%s got %s error: %s", *bucket.Name, code, message)
				continue
			}

			// generic error
			fmt.Println(*bucket.Name, err)
			continue
		}
		if out == nil {
			fmt.Println(*bucket.Name, "- does not have public access block configuration")
			continue
		}

		conf := out.PublicAccessBlockConfiguration
		if !conf.BlockPublicAcls {
			fmt.Println(*bucket.Name, "- does not block public ACLs")
		}
		if !conf.BlockPublicPolicy {
			fmt.Println(*bucket.Name, "- does not block public policy")
		}
		if !conf.IgnorePublicAcls {
			fmt.Println(*bucket.Name, "- does not ignore public ACLs")
		}
		if !conf.RestrictPublicBuckets {
			fmt.Println(*bucket.Name, "- does not restrict public buckets")
		}
	}
	return nil
}
