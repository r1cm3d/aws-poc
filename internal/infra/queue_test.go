package infra

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"testing"
)

func setup() *Queue {
	env := loadConf()
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

func TestSendMessageIntegration(t *testing.T) {
	skipShort(t)
	q := setup()
	attrs := map[string]interface{}{
		"FirstAttribute":  "Some string",
		"SecondAttribute": 666,
	}

	if _, err := q.SendMessage("Message Body", MessageAttributes(attrs)); err != nil {
		log.Fatal(err)
	}

	log.Print("successed!")
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}
