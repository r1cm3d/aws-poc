package awsstorage

import (
	"aws-poc/internal/protocol"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
)

type attstorage struct {
	session client.ConfigProvider
}

func (a attstorage) list(cid string, bucket string, path string) ([]protocol.File, error) {
	svc := s3.New(a.session)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		return nil, err
	}

	var files []protocol.File
	for _, item := range resp.Contents {
		files = append(files, protocol.File{Key: *item.Key})
	}

	return files, nil
}

func (a attstorage) get(cid string, bucket string, key string) (*protocol.File, error) {
	file, err := os.Create(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(a.session)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	size := stat.Size()
	buffer := make([]byte, size)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return &protocol.File{Key: file.Name(), Bytes: buffer}, nil
}
