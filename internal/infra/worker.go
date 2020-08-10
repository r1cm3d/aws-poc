package infra

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type worker struct {
	handler  handler
	consumer consumer
	poller   poller
}
type handler interface {
	handleMessage(msg *sqs.Message) error
}
type consumer interface {
	receiveMessage(opts ...receiveMessageInput) ([]*sqs.Message, error)
	deleteMessage(receiptHandle *string) error
}
type poller interface {
	run(consumer, handler, []*sqs.Message)
}

func Start(w worker, errors, success chan int) {
	for {
		var e, s int
		messages, err := w.consumer.receiveMessage(maxNumberOfMessages(10))
		if err != nil {
			log.Println(err)
			e++
			errors <- e
			continue
		}
		if len(messages) > 0 {
			w.poller.run(w.consumer, w.handler, messages)
			s++
			success <- s
		}
	}
}
