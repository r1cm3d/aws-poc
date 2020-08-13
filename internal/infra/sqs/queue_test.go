package sqs

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"strconv"
	"testing"
)

const (
	failMsg   = "an error should be return"
	queueName = "testQueue"
)

var errMock = errors.New("mocked error")

type sqsMock struct{}

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

func newSession() *sqs.SQS {
	env, _ := loadConf("../../../scripts/env/")
	return sqs.New(session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(env["REGION"]),
		Endpoint: aws.String(env["ENDPOINT"]),
	})))
}

func setup() *Queue {
	s := newSession()
	_, err := s.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})
	q, err := New(s, queueName)
	if err != nil {
		log.Fatal(err)
	}

	return q
}

func teardown() {
	s := newSession()
	url, _ := getQueueURL(s, queueName)
	_, err := s.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: url,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func mockQueue() *Queue {
	return &Queue{
		url: nil,
		sqs: sqsMock{},
	}
}

func TestChangeMessageVisibility(t *testing.T) {
	receiptHandle, visibilityTimeout, q := "receiptHandle", int64(10), mockQueue()

	if err := q.ChangeMessageVisibility(&receiptHandle, visibilityTimeout); err != nil {
		assert.Fail(t, failMsg)
	}
}

func TestMessageAttributeValue_Panic(t *testing.T) {
	panicFunc := func() {
		MessageAttributeValue(rune(1))
	}

	assert.Panics(t, panicFunc, failMsg)
}

func TestMessageAttributeValue(t *testing.T) {
	s, i64, i := "string", int64(666), 42
	m := map[interface{}]interface{}{
		s: &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeString),
			StringValue: aws.String(s),
		},
		i64: &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeNumber),
			StringValue: aws.String(strconv.FormatInt(i64, 10)),
		},
		i: &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeNumber),
			StringValue: aws.String(strconv.FormatInt(int64(i), 10)),
		},
	}
	b := &sqs.MessageAttributeValue{
		DataType:    aws.String(dataTypeBinary),
		BinaryValue: []byte{0},
	}

	for in, out := range m {
		act := MessageAttributeValue(in)
		assert.True(t, reflect.DeepEqual(act, out))
	}
	assert.True(t, reflect.DeepEqual(b, MessageAttributeValue([]byte{0})))
}

func TestNewWithError(t *testing.T) {
	_, err := New(sqsMock{}, "queueName")
	if err == nil {
		assert.Fail(t, failMsg)
	}
}

func TestMaxNumberOfMessages(t *testing.T) {
	exp := int64(1)
	f := maxNumberOfMessages(exp)
	in := &sqs.ReceiveMessageInput{}

	f(in)

	assert.Equal(t, aws.Int64(exp), in.MaxNumberOfMessages)
}

func TestDeleteMessage(t *testing.T) {
	q := mockQueue()
	s := "aReceiptHandle"

	err := q.deleteMessage(&s)

	if err != nil {
		assert.Fail(t, failMsg)
	}
}

func TestReceiveMessageError(t *testing.T) {
	f := func(c *sqs.ReceiveMessageInput) {}
	q := mockQueue()

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
	defer teardown()
	sendMsg(q, "Message Body")

	if msg, _ := q.receiveMessage(); msg != nil {
		for _, msg := range msg {
			assert.NotZero(t, msg.Body)
		}
	}
}

func TestSendMessageIntegration(t *testing.T) {
	skipShort(t)
	q := setup()
	defer teardown()

	sendMsg(q, "Payload")

	t.Log("succeeded!")
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
