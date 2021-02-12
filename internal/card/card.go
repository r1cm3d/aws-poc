package card

import (
	"aws-poc/internal/protocol"
)

type (
	Repository interface {
		Get(dispute *protocol.Dispute) (*protocol.Card, error)
	}
)
