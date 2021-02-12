package card

import (
	"aws-poc/internal/protocol"
)

type (
	Service interface {
		Get(dispute *protocol.Dispute) (*protocol.Card, error)
	}
)
