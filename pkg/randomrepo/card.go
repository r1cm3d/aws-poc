package randomrepo

import (
	"aws-poc/internal/protocol"
	"math/rand"
)

type randomRepository struct{}

const (
	panLength = 16
	digits    = "0123456789"
)

func (r randomRepository) Get(dispute *protocol.Dispute) (*protocol.Card, error) {
	b := make([]byte, panLength)
	for i := range b {
		b[i] = digits[rand.Int63()%int64(len(digits))]
	}
	return &protocol.Card{Number: string(b)}, nil
}
