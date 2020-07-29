package infra

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"strconv"
)

type sqsAdapter interface {
	GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error)
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	ChangeMessageVisibility(*sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error)
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

// A Queue is an sqs queue which holds queue url in url.
// Queue allows you to call actions without queue url for every call.
type Queue struct {
	url *string
	sqs sqsAdapter
}

// The DataType is a type of data used in Attributes and Message Attributes.
const (
	DataTypeString = "String"
	DataTypeNumber = "Number"
	DataTypeBinary = "Binary"
)

// New initializes Queue with queue name name.
func New(s sqsiface.SQSAPI, name string) (*Queue, error) {
	u, err := getQueueURL(s, name)
	if err != nil {
		return nil, err
	}

	return &Queue{
		sqs: s,
		url: u,
	}, nil
}

// SendMessageInput type is an adapter to change a parameter in
// sqs.SendMessageInput.
type SendMessageInput func(req *sqs.SendMessageInput)

// ChangeMessageVisibility changes a message visibility timeout.
func (q *Queue) ChangeMessageVisibility(receiptHandle *string, visibilityTimeout int64) error {
	req := &sqs.ChangeMessageVisibilityInput{
		ReceiptHandle:     receiptHandle,
		VisibilityTimeout: aws.Int64(visibilityTimeout),
		QueueUrl:          q.url,
	}
	_, err := q.sqs.ChangeMessageVisibility(req)
	return err
}

// MessageAttributes returns a SendMessageInput that changes MessageAttributes to attrs.
// A string value in attrs sets to DataTypeString.
// A []byte value in attrs sets to DataTypeBinary.
// A int and int64 value in attrs sets to DataTypeNumber. Other types cause panicking.
func MessageAttributes(attrs map[string]interface{}) SendMessageInput {
	return func(req *sqs.SendMessageInput) {
		if len(attrs) == 0 {
			return
		}

		ret := make(map[string]*sqs.MessageAttributeValue)
		for n, v := range attrs {
			ret[n] = MessageAttributeValue(v)
		}
		req.MessageAttributes = ret
	}
}

// MessageAttributeValue returns a appropriate sqs.MessageAttributeValue by type assersion of v.
// Types except string, []byte, int64 and int cause panicking.
func MessageAttributeValue(v interface{}) *sqs.MessageAttributeValue {
	switch vv := v.(type) {
	case string:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(DataTypeString),
			StringValue: aws.String(vv),
		}
	case []byte:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(DataTypeBinary),
			BinaryValue: vv,
		}
	case int64:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(DataTypeNumber),
			StringValue: aws.String(strconv.FormatInt(vv, 10)),
		}
	case int:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(DataTypeNumber),
			StringValue: aws.String(strconv.FormatInt(int64(vv), 10)),
		}
	default:
		panic("sqs: unsupported type")
	}
}

// SendMessage sends a message to sqs queue. opts are used to change parameters for a message.
func (q *Queue) SendMessage(body string, opts ...SendMessageInput) (*sqs.SendMessageOutput, error) {
	req := &sqs.SendMessageInput{
		MessageBody: aws.String(body),
		QueueUrl:    q.url,
	}

	for _, f := range opts {
		f(req)
	}

	return q.sqs.SendMessage(req)
}

type receiveMessageInput func(req *sqs.ReceiveMessageInput)

func maxNumberOfMessages(n int64) receiveMessageInput {
	return func(req *sqs.ReceiveMessageInput) {
		req.MaxNumberOfMessages = aws.Int64(n)
	}
}

func (q *Queue) receiveMessage(opts ...receiveMessageInput) ([]*sqs.Message, error) {
	req := &sqs.ReceiveMessageInput{
		QueueUrl: q.url,
	}

	for _, f := range opts {
		f(req)
	}

	resp, err := q.sqs.ReceiveMessage(req)
	if err != nil {
		return nil, err
	}
	return resp.Messages, nil
}

func (q *Queue) deleteMessage(receiptHandle *string) error {
	_, err := q.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      q.url,
		ReceiptHandle: receiptHandle,
	})
	return err
}

func getQueueURL(s sqsAdapter, name string) (*string, error) {
	req := &sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	}

	resp, err := s.GetQueueUrl(req)
	if err != nil {
		return nil, err
	}
	return resp.QueueUrl, nil
}
