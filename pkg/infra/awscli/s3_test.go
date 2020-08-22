package awscli

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"testing"
)

const (
	bucketName = "bucket"
	key        = "kellyKey"
	filename   = "../../../scripts/env/.env"
)

func TestUploadIntegration(t *testing.T) {
	skipShort(t)
	setupBucket()
	defer cleanupBucket()

	sess := newLocalSessionWithS3ForcePathStyle()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("enable to open file")
	}
	defer file.Close()

	s3cli := S3cli{
		sess,
	}

	err = s3cli.Upload(bucketName, key, file)
	if err != nil {
		t.Errorf("error on Upload = %v", err)
	}
}

func TestListIntegration(t *testing.T) {
	skipShort(t)
	setupBucket()
	defer cleanupBucket()

	sess := newLocalSessionWithS3ForcePathStyle()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("enable to open file")
	}
	defer file.Close()

	s3cli := S3cli{
		sess,
	}

	err = s3cli.Upload(bucketName, key, file)
	if err != nil {
		t.Errorf("error on Upload = %v", err)
	}

	err = s3cli.List(bucketName, key)
	if err != nil {
		t.Errorf("error on List = %v", err)
	}
}

func TestGetIntegration(t *testing.T) {
	skipShort(t)
	setupBucket()
	defer cleanupBucket()

	file, err := os.Open("../../../scripts/env/.env")
	if err != nil {
		log.Fatal("enable to open file")
	}
	defer file.Close()

	s3cli := S3cli{
		newLocalSessionWithS3ForcePathStyle(),
	}

	err = s3cli.Upload(bucketName, key, file)
	if err != nil {
		t.Errorf("error on Upload = %v", err)
	}

	err = s3cli.Get(bucketName, key)
	if err != nil {
		t.Errorf("error on get = %v", err)
	}
}

func setupBucket() {
	sess := newLocalSessionWithS3ForcePathStyle()

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
	sess := newLocalSessionWithS3ForcePathStyle()

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

	_ = os.Remove(key)
}
