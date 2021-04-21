package awsstorage

import (
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestAttListIntegration(t *testing.T) {
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

	attsto := attstorage{
		sess,
	}

	if err := s3cli.Upload(bucketName, key, file); err != nil {
		t.Errorf("error on Upload = %v", err)
	}
	if files, err := attsto.list("cid", bucketName, key); err != nil {
		t.Errorf("error on List = %v", err)
	} else {
		for i, f := range files {
			fmt.Printf("File %d:%v retrieved", i, f.Key)
		}
	}
}

func TestAttGetIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupBucket()
	defer cleanupBucket()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("enable to open attachment")
	}
	defer file.Close()

	sess := awssession.NewLocalSessionWithS3ForcePathStyle()
	s3cli := S3cli{
		sess,
	}
	attsto := attstorage{
		sess,
	}

	if err := s3cli.Upload(bucketName, key, file); err != nil {
		t.Errorf("error on Upload = %v", err)
	}
	if file, err := attsto.get("cid", bucketName, key); err != nil {
		t.Errorf("error on get = %v", err)
	} else {
		fmt.Printf("File %s retrieved. Bytes: %v", file.Key, file.Bytes)
	}
}
