package network

import "aws-poc/internal/protocol"

type Creator interface {
	Create(*protocol.Dispute, *protocol.Card, *protocol.Attachment) (*protocol.Chargeback, error)
}
