package infra

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleMessage_Error(t *testing.T) {
	cp := concurPoller{}

	act := cp.handleMessage(nil, errHandler{}, nil)

	assert.Equal(t, err, act)
}

func TestHandleMessage(t *testing.T) {
	cp := concurPoller{}

	act := cp.handleMessage(fakeErrQueue{}, okHandler{}, &sqs.Message{ReceiptHandle: aws.String("receipt")})

	assert.True(t, act == nil)
}

func TestRun(t *testing.T) {
	cp := concurPoller{}

	cp.run(fakeOkQueue{}, okHandler{}, []*sqs.Message{{ReceiptHandle: aws.String("receipt")}})
}

func TestRun_Error(t *testing.T) {
	cp := concurPoller{}

	cp.run(fakeOkQueue{}, errHandler{}, []*sqs.Message{{ReceiptHandle: aws.String("receipt")}})
}
