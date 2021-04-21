package awsmessaging

import (
	"log"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type (
	worker struct {
		handler
		consumer
		poller
	}

	handler interface {
		handleMessage(cid, body string) error
	}

	consumer interface {
		receiveMessage(opts ...receiveMessageInput) ([]*sqs.Message, error)
		deleteMessage(receiptHandle *string) error
	}

	poller interface {
		run(consumer, handler, []*sqs.Message)
	}
)

// Start starts worker passed as argument and errors and success through channels
func Start(w worker, errors, success chan<- int) {
	for {
		var e, s int
		messages, err := w.receiveMessage(maxNumberOfMessages(10))
		if err != nil {
			log.Println(err)
			e++
			errors <- e
			continue
		}
		if len(messages) > 0 {
			w.run(w, w, messages)
			s++
			success <- s
		}
	}
}
