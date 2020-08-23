package awscli

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type errHandler func(*sqs.Message) error
type okHandler func(*sqs.Message) error
type fakeErrQueue struct{}
type fakeOkQueue struct{}
type fakePoller struct{}

var err = errors.New("sambarilove")

func (m errHandler) handleMessage(_ *sqs.Message) error {
	return err
}

func (m okHandler) handleMessage(_ *sqs.Message) error {
	return nil
}

func (q fakeErrQueue) deleteMessage(_ *string) error {
	return nil
}

func (q fakeErrQueue) receiveMessage(_ ...receiveMessageInput) ([]*sqs.Message, error) {
	return nil, err
}

func (q fakeOkQueue) deleteMessage(_ *string) error {
	return nil
}

func (q fakeOkQueue) receiveMessage(_ ...receiveMessageInput) ([]*sqs.Message, error) {
	return []*sqs.Message{{Body: aws.String("")}}, nil
}

func (p fakePoller) run(consumer, handler, []*sqs.Message) {}

