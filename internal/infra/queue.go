package infra

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// A Queue is an SQS queue which holds queue url in URL.
// Queue allows you to call actions without queue url for every call.
type Queue struct {
	SQS sqsiface.SQSAPI
	URL *string
}

// New initializes Queue with queue name name.
func New(s sqsiface.SQSAPI, name string) (*Queue, error) {
	u, err := getQueueURL(s, name)
	if err != nil {
		return nil, err
	}

	return &Queue{
		SQS: s,
		URL: u,
	}, nil
}

// The SendMessageInput type is an adapter to change a parameter in
// sqs.SendMessageInput.
type SendMessageInput func(req *sqs.SendMessageInput)

// ChangeMessageVisibility changes a message visibility timeout.
func (q *Queue) ChangeMessageVisibility(receiptHandle *string, visibilityTimeout int64) error {
	req := &sqs.ChangeMessageVisibilityInput{
		ReceiptHandle:     receiptHandle,
		VisibilityTimeout: aws.Int64(visibilityTimeout),
		QueueUrl:          q.URL,
	}
	_, err := q.SQS.ChangeMessageVisibility(req)
	return err
}

// SendMessage sends a message to SQS queue. opts are used to change parameters for a message.
func (q *Queue) SendMessage(body string, opts ...SendMessageInput) (*sqs.SendMessageOutput, error) {
	req := &sqs.SendMessageInput{
		MessageBody: aws.String(body),
		QueueUrl:    q.URL,
	}

	for _, f := range opts {
		f(req)
	}

	return q.SQS.SendMessage(req)
}

type receiveMessageInput func(req *sqs.ReceiveMessageInput)

func maxNumberOfMessages(n int64) receiveMessageInput {
	return func(req *sqs.ReceiveMessageInput) {
		req.MaxNumberOfMessages = aws.Int64(n)
	}
}

func (q *Queue) receiveMessage(opts ...receiveMessageInput) ([]*sqs.Message, error) {
	req := &sqs.ReceiveMessageInput{
		QueueUrl: q.URL,
	}

	for _, f := range opts {
		f(req)
	}

	resp, err := q.SQS.ReceiveMessage(req)
	if err != nil {
		return nil, err
	}
	return resp.Messages, nil
}

func (q *Queue) deleteMessage(receiptHandle *string) error {
	_, err := q.SQS.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      q.URL,
		ReceiptHandle: receiptHandle,
	})
	return err
}

func getQueueURL(s sqsiface.SQSAPI, name string) (*string, error) {
	req := &sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	}

	resp, err := s.GetQueueUrl(req)
	if err != nil {
		return nil, err
	}
	return resp.QueueUrl, nil
}
