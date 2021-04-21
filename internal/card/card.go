package card

import (
	"aws-poc/internal/protocol"
)

// A Service provides interactions with protocol.Card.
type Service interface {
	Get(*protocol.Dispute) (*protocol.Card, error)
}
