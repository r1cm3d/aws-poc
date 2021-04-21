package network

import "aws-poc/internal/protocol"

// A Creator is responsible for create a protocol.Chargeback in the network brand given a protocol.Dispute,
// protocol.Card and protocol.Attachment.
type Creator interface {
	Create(*protocol.Dispute, *protocol.Card, *protocol.Attachment) (*protocol.Chargeback, error)
}
