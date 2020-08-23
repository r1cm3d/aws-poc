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
type sqsMock struct{}

var errMock = errors.New("mocked error")

func (s sqsMock) GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return nil, errMock
}
func (s sqsMock) ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return nil, errMock
}
func (s sqsMock) ChangeMessageVisibility(*sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error) {
	return nil, nil
}
func (s sqsMock) SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return nil, nil
}
func (s sqsMock) DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return nil, nil
}

func (m errHandler) handleMessage(_ *sqs.Message) error {
	return errMock
}

func (m okHandler) handleMessage(_ *sqs.Message) error {
	return nil
}

func (q fakeErrQueue) deleteMessage(_ *string) error {
	return nil
}

func (q fakeErrQueue) receiveMessage(_ ...receiveMessageInput) ([]*sqs.Message, error) {
	return nil, errMock
}

func (q fakeOkQueue) deleteMessage(_ *string) error {
	return nil
}

func (q fakeOkQueue) receiveMessage(_ ...receiveMessageInput) ([]*sqs.Message, error) {
	return []*sqs.Message{{Body: aws.String("")}}, nil
}

func (p fakePoller) run(consumer, handler, []*sqs.Message) {}
