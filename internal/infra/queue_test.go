package infra

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const failMsg = "an error should be return"

var mockedError = errors.New("mocked error")

func setup() *Queue {
	env, _ := loadConf("../../scripts/env/")
	s := sqs.New(session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
	})))
	q, err := New(s, env["SQS_TEST_QUEUE"])
	if err != nil {
		log.Fatal(err)
	}

	return q
}

type sqsMock struct{}

func (s sqsMock) GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return nil, mockedError
}
func (s sqsMock) ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return nil, mockedError
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

func TestNewWithError(t *testing.T) {
	_, err := New(sqsMock{}, "queueName")
	if err == nil {
		assert.Fail(t, failMsg)
	}
}

func TestSendMessageIntegration(t *testing.T) {
	skipShort(t)
	q := setup()

	sendMsg(q, "Payload")

	log.Print("successed!")
}

func TestMaxNumberOfMessages(t *testing.T) {
	exp := int64(1)
	f := maxNumberOfMessages(exp)
	in := &sqs.ReceiveMessageInput{}

	f(in)

	assert.Equal(t, aws.Int64(exp), in.MaxNumberOfMessages)
}

func TestReceiveMessageError(t *testing.T) {
	f := func (c *sqs.ReceiveMessageInput){

	}
	q := &Queue{
		url: nil,
		sqs: sqsMock{},
	}

	_, err := q.receiveMessage(f)
	if err == nil {
		assert.Fail(t, failMsg)
	}
}

func TestMessageAttributes(t *testing.T) {
	f := MessageAttributes(map[string]interface{}{})
	f(nil)
}

func TestReceiveMessageIntegration(t *testing.T) {
	skipShort(t)
	q := setup()
	sendMsg(q, "Message Body")

	if msgs, _ := q.receiveMessage(); msgs != nil {
		for _, msg := range msgs {
			assert.NotZero(t, msg.Body)
		}
	}
}

func sendMsg(queue *Queue, body string) {
	attrs := map[string]interface{}{
		"FirstAttribute":  "Some string",
		"SecondAttribute": 666,
	}
	if _, err := queue.SendMessage(body, MessageAttributes(attrs)); err != nil {
		log.Fatal(err)
	}
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}
