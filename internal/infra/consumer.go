package infra

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"sync"
)

// Handle is responsible to deal with message
type Handle func(msg *sqs.Message) error
type consumer interface {
	receiveMessage(opts ...receiveMessageInput) ([]*sqs.Message, error)
	deleteMessage(receiptHandle *string) error
}

// HandleMessage delegates message handling to Handle
func (f Handle) HandleMessage(msg *sqs.Message) error {
	return f(msg)
}

// Handler is an abstraction that handle sqs messages
type Handler interface {
	HandleMessage(msg *sqs.Message) error
}

// Start starts a worker giving a Queue and a Handler
func Start(c consumer, h Handler, done chan bool) {
	for {
		messages, err := c.receiveMessage(maxNumberOfMessages(10))
		if err != nil {
			log.Println(err)
			done <- true
			continue
		}
		if len(messages) > 0 {
			run(c, h, messages)
		}
	}
}

func run(c consumer, h Handler, messages []*sqs.Message) {
	numMessages := len(messages)

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			defer wg.Done()
			if err := handleMessage(c, m, h); err != nil {
				log.Println(err)
			}
		}(messages[i])
	}

	wg.Wait()
}

func handleMessage(c consumer, m *sqs.Message, h Handler) error {
	var err error
	err = h.HandleMessage(m)
	if err != nil {
		return err
	}
	return c.deleteMessage(m.ReceiptHandle)
}
