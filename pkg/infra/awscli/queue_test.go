package awscli

import (
	"github.com/aws/aws-sdk-go/aws"
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

func newSQS() *sqs.SQS {
	return sqs.New(newLocalSession())
}

func setup() *Queue {
	s := newSQS()
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	}
	if _, err := s.CreateQueue(input); err != nil {
		panic(err)
	}
	q, err := New(s, queueName)
	if err != nil {
		log.Fatal(err)
	}

	return q
}

func teardown() {
	s := newSQS()
	url, _ := getQueueURL(s, queueName)
	input := &sqs.DeleteQueueInput{
		QueueUrl: url,
	}
	if _, err := s.DeleteQueue(input); err != nil {
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
		t.Error(failMsg)
	}
}

func TestMessageAttributeValue_Panic(t *testing.T) {
	panicFunc := func() {
		MessageAttributeValue(rune(1))
	}

	if !panic(panicFunc) {
		t.Error("function should panic")
	}
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
	if _, err := New(sqsMock{}, "queueName"); err == nil {
		t.Error(failMsg)
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

	if err := q.deleteMessage(&s); err != nil {
		assert.Fail(t, failMsg)
	}
}

func TestReceiveMessageError(t *testing.T) {
	f := func(c *sqs.ReceiveMessageInput) {}
	q := mockQueue()

	if _, err := q.receiveMessage(f); err == nil {
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

func panic(f func()) (ok bool) {
	ok = false
	defer func() {
		if r := recover(); r != nil {
			ok = true
		}
	}()
	f()

	return
}
