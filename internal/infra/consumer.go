package infra

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"sync"
)

type Handle func(msg *sqs.Message) error

func (f Handle) HandleMessage(msg *sqs.Message) error {
	return f(msg)
}

type Handler interface {
	HandleMessage(msg *sqs.Message) error
}

func Start(q *Queue, h Handler) {
	for {
		messages, err := q.ReceiveMessage(MaxNumberOfMessages(10))
		if err != nil {
			log.Println(err)
			continue
		}
		if len(messages) > 0 {
			run(q, h, messages)
		}
	}
}

func run(q *Queue, h Handler, messages []*sqs.Message) {
	numMessages := len(messages)

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			defer wg.Done()
			if err := handleMessage(q, m, h); err != nil {
				log.Println(err)
			}
		}(messages[i])
	}

	wg.Wait()
}

func handleMessage(q *Queue, m *sqs.Message, h Handler) error {
	var err error
	err = h.HandleMessage(m)
	if err != nil {
		return err
	}
	return q.DeleteMessage(m.ReceiptHandle)
}
