package infra

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type errHandler struct{}
type okHandler struct{}
type fakeQueue struct{}

var err = errors.New("sambarilove")

func (m errHandler) HandleMessage(_ *sqs.Message) error {
	return err
}

func (m okHandler) HandleMessage(_ *sqs.Message) error {
	return nil
}

func (q fakeQueue) deleteMessage(_ *string) error {
	return nil
}

func (q fakeQueue) receiveMessage(_ ...receiveMessageInput) ([]*sqs.Message, error) {
	return nil, err
}

func TestHandleMessage_Error(t *testing.T) {
	h := errHandler{}

	act := handleMessage(nil, nil, h)

	assert.Equal(t, err, act)
}

func TestHandleMessage(t *testing.T) {
	q, h, m := &fakeQueue{}, okHandler{}, &sqs.Message{ReceiptHandle: aws.String("receipt")}

	act := handleMessage(q, m, h)

	assert.True(t, act == nil)
}

func TestStart(t *testing.T) {
	q, h, d := &fakeQueue{}, okHandler{}, make(chan bool)

	go Start(q, h, d)

	assert.True(t, <-d)
}
