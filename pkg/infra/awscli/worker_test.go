package awscli

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestStart_Error(t *testing.T) {
	e, s, w := make(chan int), make(chan int), worker{
		consumer: fakeErrQueue{},
	}

	go Start(w, e, s)

	assert.Equal(t, <-e, 1)
}

func TestStart(t *testing.T) {
	e, s, w := make(chan int), make(chan int), worker{
		consumer: fakeOkQueue{},
		poller:   fakePoller{},
	}

	go Start(w, e, s)

	assert.Equal(t, <-s, 1)
}
