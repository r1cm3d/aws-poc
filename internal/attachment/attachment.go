package attachment

import (
	"aws-poc/internal/protocol"
	"fmt"
)

const filenameRoot = "disputes"

type (
	file struct {
		key string
	}
	Service interface {
		Get(dispute *protocol.Dispute) (*protocol.Attachment, error)
		Save(chargeback *protocol.Chargeback) error
	}

	storage interface {
		list(cid string, bucket string, path string) ([]file, error)
		get(cid string, bucket string, key string) (*file, error)
	}

	repository interface {
		getUnsentFiles(*protocol.Dispute, []file) ([]file, error)
	}

	archiver interface {
		compact(cid string, files []file, strToRemove string) (*protocol.Attachment, error)
	}

	svc struct {
		storage
		archiver
		repository
	}
)

func (s svc) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	var (
		files          []file
		err            error
		unsentFiles    []file
		filesToCompact []file
	)
	path := fmt.Sprintf("%s/%d/%d", filenameRoot, dispute.AccountId, dispute.DisputeId)
	if files, err = s.list(dispute.Cid, dispute.OrgId, path); err != nil {
		return nil, err
	}
	if unsentFiles, err = s.getUnsentFiles(dispute, files); err != nil {
		return nil, err
	}

	for _, uf := range unsentFiles {
		var rf *file
		if rf, err = s.get(dispute.Cid, dispute.OrgId, uf.key); err != nil {
			return nil, err
		}
		filesToCompact = append(filesToCompact, *rf)
	}

	return s.compact(dispute.Cid, filesToCompact, path)
}

func (s svc) Save(chargeback *protocol.Chargeback) error {
	// TODO: implement it
	return nil
}
