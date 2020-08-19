package awscli

import (
	"aws-poc/internal/infra"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"testing"
)

const (
	bucketName = "bucket"
	key = "kellyKey"
)

func TestUploadIntegration(t *testing.T) {
	skipShort(t)
	setupBucket()
	defer cleanupBucket()

	env, _ := infra.LoadDefaultConf()
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
		S3ForcePathStyle: aws.Bool(true),
	}))

	file, err := os.Open("../../../scripts/env/.env")
	if err != nil {
		log.Fatal("enable to open file")
	}

	defer file.Close()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(key),
		Body: file,
	})
	if err != nil {
		log.Fatal(err)
	}
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

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(key)})
	if err != nil {
		log.Fatal("unable to delete object from bucket")
	}
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
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
