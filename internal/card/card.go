package card

import (
	"aws-poc/internal/protocol"
)

type Service interface {
		Get(*protocol.Dispute) (*protocol.Card, error)
}

