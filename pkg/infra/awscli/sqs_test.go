package awscli

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type errHandler func(string) error
type okHandler func(string) error
type fakeErrQueue struct{}
type fakeOkQueue struct{}
type fakePoller struct{}
type sqsMock struct{}

var (
	errMock       = errors.New("mocked error")
	attr          = map[string]*sqs.MessageAttributeValue{"cid": {StringValue: aws.String("2ce488cd-e6b0-4fea-a960-31256018cf08")}}
	mockedMessage = &sqs.Message{MessageAttributes: attr, Body: aws.String(""), ReceiptHandle: aws.String("receipt")}
)

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

func (m errHandler) handleMessage(_, _ string) error {
	return errMock
}

func (m okHandler) handleMessage(_, _ string) error {
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
