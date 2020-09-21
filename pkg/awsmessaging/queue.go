package awsmessaging

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	dataTypeString = "String"
	dataTypeNumber = "Number"
	dataTypeBinary = "Binary"
)

type (
	// A Queue is an sqs queue which holds queue url in url.
	// Queue allows you to call actions without queue url for every call.
	Queue struct {
		url *string
		sqs sqsAdapter
	}

	// SendMessageInput type is an adapter to change a parameter in
	// sqs.SendMessageInput.
	SendMessageInput func(req *sqs.SendMessageInput)

	receiveMessageInput func(req *sqs.ReceiveMessageInput)

	sqsAdapter interface {
		// This name is not linter compliance because is equal to GetQueueUrl method
		GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error)
		ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
		ChangeMessageVisibility(*sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error)
		SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
		DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	}
)

// New initializes Queue with queue name.
func New(s sqsAdapter, name string) (*Queue, error) {
	u, err := getQueueURL(s, name)
	if err != nil {
		return nil, err
	}

	return &Queue{
		sqs: s,
		url: u,
	}, nil
}

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
// A string value in attrs sets to dataTypeString.
// A []byte value in attrs sets to dataTypeBinary.
// A int and int64 value in attrs sets to dataTypeNumber. Other types cause panicking.
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

// MessageAttributeValue returns a appropriate sqs.MessageAttributeValue by type assertion of v.
// Types except string, []byte, int64 and int cause panicking.
func MessageAttributeValue(v interface{}) *sqs.MessageAttributeValue {
	switch vv := v.(type) {
	case string:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeString),
			StringValue: aws.String(vv),
		}
	case []byte:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeBinary),
			BinaryValue: vv,
		}
	case int64:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeNumber),
			StringValue: aws.String(strconv.FormatInt(vv, 10)),
		}
	case int:
		return &sqs.MessageAttributeValue{
			DataType:    aws.String(dataTypeNumber),
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
