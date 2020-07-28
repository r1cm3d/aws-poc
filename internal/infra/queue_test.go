package infra

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

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

func TestSendMessageIntegration(t *testing.T) {
	skipShort(t)
	q := setup()

	sendMsg(q, "Payload")

	log.Print("successed!")
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
