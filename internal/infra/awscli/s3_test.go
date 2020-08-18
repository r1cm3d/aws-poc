package awscli

import (
	"aws-poc/internal/infra"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"testing"
)

const bucketName = "bucket"

func TestUploadIntegration(t *testing.T) {
	skipShort(t)
	setupBucket()
	fmt.Println("test bucket")
	cleanupBucket()
}

func setupBucket() {
	env, _ := infra.LoadDefaultConf()
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
		S3ForcePathStyle: aws.Bool(true),
	}))

	svc := s3.New(sess)
	_, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName)})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Fatal(err)
	}
}

func cleanupBucket() {
	env, _ := infra.LoadDefaultConf()
	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
		S3ForcePathStyle: aws.Bool(true),
	})

	svc := s3.New(sess)

	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Fatal(err)
	}
}
