package awsstorage

import (
	"aws-poc/internal/attachment"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
)

type attstorage struct {
	session client.ConfigProvider
}

func (a attstorage) list(cid string, bucket string, path string) ([]attachment.File, error) {
	// TODO: log here
	svc := s3.New(a.session)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		// TODO: log here
		return nil, err
	}

	// TODO: log here
	var files []attachment.File
	for _, item := range resp.Contents {
		// TODO: log here
		files = append(files, attachment.File{Key: *item.Key})
	}

	// TODO: log here
	return files, nil
}
