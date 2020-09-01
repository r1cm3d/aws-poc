package awscli

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"sync"
)

type concurPoller struct{}

func (p concurPoller) run(c consumer, h handler, messages []*sqs.Message) {
	numMessages := len(messages)

	var wg sync.WaitGroup

	wg.Add(numMessages)
	for i := range messages {
		go handle(p, c, h, &wg, messages[i])
	}
	wg.Wait()
}

func (p concurPoller) handleMessage(c consumer, h handler, m *sqs.Message) error {
	var cid string
	if v, ok := m.MessageAttributes["cid"]; ok {
		cid = *v.StringValue
	}
	if err := h.handleMessage(cid, *m.Body); err != nil {
		return err
	}
	return c.deleteMessage(m.ReceiptHandle)
}

func handle(p concurPoller, c consumer, h handler, wg *sync.WaitGroup, m *sqs.Message) {
	defer wg.Done()
	if err := p.handleMessage(c, h, m); err != nil {
		log.Println(err)
	}
}
