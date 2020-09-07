package awscli

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
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

	if !didPanic(panicFunc) {
		t.Error("function should panic")
	}
}

func TestMessageAttributeValue(t *testing.T) {
	s, i64, integer := "string", int64(666), 42
	m := map[interface{}]interface{}{
		s: &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeString),
			StringValue: aws.String(s),
		},
		i64: &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeNumber),
			StringValue: aws.String(strconv.FormatInt(i64, 10)),
		},
		integer: &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeNumber),
			StringValue: aws.String(strconv.FormatInt(int64(integer), 10)),
		},
	}
	b := &sqs.MessageAttributeValue{
		DataType:    aws.String(dataTypeBinary),
		BinaryValue: []byte{0},
	}

	for i, want := range m {
		t.Run(fmt.Sprintf("%T", m[i]), func(t *testing.T) {
			got := MessageAttributeValue(i)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("want: %d; got: %d", want, got)
			}
		})
	}

	got := MessageAttributeValue([]byte{0})
	if !reflect.DeepEqual(got, b) {
		t.Errorf("want: %d; got: %d", b, got)
	}
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

	got := *in.MaxNumberOfMessages

	if exp != got {
		t.Errorf("want: %d; got: %d", exp, got)
	}
}

func TestDeleteMessage(t *testing.T) {
	q := mockQueue()
	s := "aReceiptHandle"

	if err := q.deleteMessage(&s); err != nil {
		t.Error(failMsg)
	}
}

func TestReceiveMessageError(t *testing.T) {
	f := func(c *sqs.ReceiveMessageInput) {}
	q := mockQueue()

	if _, err := q.receiveMessage(f); err == nil {
		t.Error(failMsg)
	}
}

func TestMessageAttributes(t *testing.T) {
	f := MessageAttributes(map[string]interface{}{})
	f(nil)
}

func TestReceiveMessageIntegration(t *testing.T) {
	skipShort(t)
	q, body := setup(), "Message Body"
	defer teardown()
	sendMsg(q, body)

	if messages, _ := q.receiveMessage(); messages != nil {
		for _, msg := range messages {
			got := *msg.Body
			if body != got {
				t.Errorf("want: %s; got: %s", body, got)
			}
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

func didPanic(f func()) (ok bool) {
	ok = false
	defer func() {
		if r := recover(); r != nil {
			ok = true
		}
	}()
	f()

	return
}
