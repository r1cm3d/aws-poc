package attachment

import "aws-poc/internal/protocol"

type (
	Service interface {
		Get(dispute *protocol.Dispute) (*protocol.Attachment, error)
		Save(chargeback *protocol.Chargeback) error
	}
)
