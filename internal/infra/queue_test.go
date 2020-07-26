package infra

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"testing"
)

func TestSendMessageIntegration(t *testing.T) {
	s := sqs.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})))
	q, err := New(s, "test-queue")
	if err != nil {
		log.Fatal(err)
	}
	attrs := map[string]interface{}{
		"ATTR1": "STRING!!",
		"ATTR2": 12345,
	}

	if _, err := q.SendMessage("MESSAGE BODY", MessageAttributes(attrs)); err != nil {
		log.Fatal(err)
	}

	log.Print("successed!")
}
