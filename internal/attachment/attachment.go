package attachment

import "aws-poc/internal/protocol"

type (
	Repository interface {
		Get(dispute *protocol.Dispute) (*protocol.Attachment, error)
		Save(chargeback *protocol.Chargeback) error
	}
)
