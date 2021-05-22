package awsrepository

import (
	"aws-poc/internal/protocol"
	"fmt"
)

type (
	attachmentRepository struct {
		repository
	}
)

func (r attachmentRepository) Get(dispute *protocol.Dispute, allFiles []protocol.File) ([]protocol.Attachment, error) {
	var attachments []protocol.Attachment
	for _, f := range allFiles {
		att, err := r.query("ID", r.ID(*dispute, f.Key), protocol.Attachment{})

		if err != nil {
			return []protocol.Attachment{}, nil
		}

		attachments = append(attachments, att.(protocol.Attachment))
	}
	return attachments, nil
}

func (r attachmentRepository) ID(dispute protocol.Dispute, key string) string {
	return fmt.Sprintf("%s::%d::%d::%s", dispute.OrgID, dispute.AccountID, dispute.DisputeID, key)
}
