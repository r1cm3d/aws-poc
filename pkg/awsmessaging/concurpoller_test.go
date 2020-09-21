package awsmessaging

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func TestHandleMessage_Error(t *testing.T) {
	cp := concurPoller{}

	if got := cp.handleMessage(nil, errHandler(nil), mockedMessage); got != errMock {
		t.Errorf("want: %d; got: %d", errMock, got)
	}
}

func TestHandleMessage(t *testing.T) {
	cp, msg := concurPoller{}, mockedMessage

	if err := cp.handleMessage(fakeErrQueue{}, okHandler(nil), msg); err != nil {
		t.Error("should not return an error at handleMessage")
	}
}

func TestRun(t *testing.T) {
	cp := concurPoller{}

	cp.run(fakeOkQueue{}, okHandler(nil), []*sqs.Message{mockedMessage})
}

func TestRun_Error(t *testing.T) {
	cp := concurPoller{}

	cp.run(fakeOkQueue{}, errHandler(nil), []*sqs.Message{mockedMessage})
}
