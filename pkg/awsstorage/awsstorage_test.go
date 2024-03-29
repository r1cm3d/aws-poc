package awsstorage

import (
	"aws-poc/internal/protocol"
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestUploadIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupBucket()
	defer cleanupBucket()

	sess := awssession.NewLocalSessionWithS3ForcePathStyle()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("enable to open attachment")
	}
	defer file.Close()

	s3cli := S3cli{
		sess,
	}

	if err := s3cli.Upload(bucketName, key, file); err != nil {
		t.Errorf("error on Upload = %v", err)
	}
}

func TestListIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupBucket()
	defer cleanupBucket()

	sess := awssession.NewLocalSessionWithS3ForcePathStyle()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("enable to open attachment")
	}
	defer file.Close()

	s3cli := S3cli{
		sess,
	}
	expFiles := [2]protocol.File{{Key: "file1"}, {Key: "file2"}}

	for _, f := range expFiles {
		if err := s3cli.Upload(bucketName, f.Key, file); err != nil {
			t.Errorf("error on Upload = %v", err)
		}
	}

	if files, err := s3cli.List("cid", bucketName, "file"); err != nil {
		t.Errorf("error on List = %v", err)
	} else if files[0].Key == expFiles[0].Key && files[1].Key == expFiles[1].Key {
		fmt.Printf("file matched %v", files)
	} else {
		t.Errorf("files do not match: %v, %v", files, expFiles)
	}
}

func TestGetIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupBucket()
	defer cleanupBucket()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("enable to open attachment")
	}
	defer file.Close()

	s3cli := S3cli{
		awssession.NewLocalSessionWithS3ForcePathStyle(),
	}

	if err := s3cli.Upload(bucketName, key, file); err != nil {
		t.Errorf("error on Upload = %v", err)
	}
	if file, err := s3cli.Get("cid", bucketName, key); err != nil || file.Key != key {
		t.Errorf("error on get = %v", err)
	}
}

func setupBucket() {
	sess := awssession.NewLocalSessionWithS3ForcePathStyle()

	svc := s3.New(sess)
	if _, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucketName)}); err != nil {
		log.Fatal(err.Error())
	}

	hbi := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	if err := svc.WaitUntilBucketExists(hbi); err != nil {
		log.Fatal(err)
	}
}

func cleanupBucket() {
	sess := awssession.NewLocalSessionWithS3ForcePathStyle()
	svc := s3.New(sess)

	files := [3]string{key, "file1", "file2"}
	for _, f := range files {
		doi := &s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(f)}
		if _, err := svc.DeleteObject(doi); err != nil {
			log.Fatal("unable to delete object from bucket")
		}
	}

	hoi := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
	if err := svc.WaitUntilObjectNotExists(hoi); err != nil {
		log.Fatal(err)
	}

	dbi := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}
	if _, err := svc.DeleteBucket(dbi); err != nil {
		log.Fatal(err)
	}

	hbi := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	if err := svc.WaitUntilBucketNotExists(hbi); err != nil {
		log.Fatal(err)
	}

	_ = os.Remove(key)
}
