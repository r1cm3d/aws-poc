package attachment

import (
	"aws-poc/internal/protocol"
	"fmt"
)

const filenameRoot = "disputes"

type (
	// A Service provides interactions with protocol.Attachments.
	Service interface {
		Get(dispute *protocol.Dispute) (*protocol.Attachment, error)
		Save(chargeback *protocol.Chargeback) error
	}

	storage interface {
		list(cid string, bucket string, path string) ([]protocol.File, error)
		Get(cid string, bucket string, key string) (*protocol.File, error)
	}

	repository interface {
		getUnsentFiles(*protocol.Dispute, []protocol.File) ([]protocol.File, error)
		save(chargeback *protocol.Chargeback) error
	}

	archiver interface {
		compact(cid string, files []protocol.File, strToRemove string) (*protocol.Attachment, error)
	}

	svc struct {
		storage
		archiver
		repository
	}
)

// NewFile creates and initiliazes a new protocol.File according key argument.
func NewFile(key string) protocol.File {
	return protocol.File{
		Key: key,
	}
}

func (s svc) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	var (
		files          []protocol.File
		err            error
		unsentFiles    []protocol.File
		filesToCompact []protocol.File
	)
	path := fmt.Sprintf("%s/%d/%d", filenameRoot, dispute.AccountID, dispute.DisputeID)
	if files, err = s.list(dispute.Cid, dispute.OrgID, path); err != nil {
		return nil, err
	}
	if unsentFiles, err = s.getUnsentFiles(dispute, files); err != nil {
		return nil, err
	}

	for _, uf := range unsentFiles {
		var rf *protocol.File
		if rf, err = s.storage.Get(dispute.Cid, dispute.OrgID, uf.Key); err != nil {
			return nil, err
		}
		filesToCompact = append(filesToCompact, *rf)
	}

	return s.compact(dispute.Cid, filesToCompact, path)
}

func (s svc) Save(chargeback *protocol.Chargeback) error {
	return s.save(chargeback)
}
