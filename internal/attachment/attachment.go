package attachment

import (
	"aws-poc/internal/protocol"
	"encoding/base64"
	"fmt"
)

const filenameRoot = "disputes"

type (
	// A Service provides interactions with protocol.Attachment
	Service interface {
		Get(dispute *protocol.Dispute) (*protocol.Attachment, error)
		Save(chargeback *protocol.Chargeback) error
	}

	// A Compressor Compress a slice of protocol.File into a []byte
	Compressor interface {
		Compress(cid string, files []protocol.File, strToRemove string) ([]byte, error)
	}

	storage interface {
		list(cid string, bucket string, path string) ([]protocol.File, error)
		Get(cid string, bucket string, key string) (*protocol.File, error)
	}

	// A Repository provides an interface to interact with attachment files
	Repository interface {
		// Get retrieves a []protocol.Attachment persisted at database given *protocol.Dispute and []protocol.File
		Get(*protocol.Dispute, []protocol.File) ([]protocol.Attachment, error)
	}

	svc struct {
		storage
		Compressor
		Repository
	}
)

// NewFile creates and initializes a new protocol.File according key argument.
func NewFile(key string) protocol.File {
	return protocol.File{
		Key: key,
	}
}

// TODO: implement it
func (s svc) getUnsentFiles(dispute *protocol.Dispute, allFiles []protocol.File, attachments []protocol.Attachment) ([]protocol.File, error) {
	return nil, nil
}

func (s svc) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	var (
		files          []protocol.File
		err            error
		attachments    []protocol.Attachment
		unsentFiles    []protocol.File
		filesToCompact []protocol.File
		compactFiles   []byte
	)
	path := fmt.Sprintf("%s/%d/%d", filenameRoot, dispute.AccountID, dispute.DisputeID)
	if files, err = s.list(dispute.Cid, dispute.OrgID, path); err != nil {
		return nil, err
	}

	// TODO: test it
	if attachments, err = s.Repository.Get(dispute, files); err != nil {
		return nil, err
	}

	if unsentFiles, err = s.getUnsentFiles(dispute, files, attachments); err != nil {
		return nil, err
	}

	for _, uf := range unsentFiles {
		var rf *protocol.File
		if rf, err = s.storage.Get(dispute.Cid, dispute.OrgID, uf.Key); err != nil {
			return nil, err
		}
		filesToCompact = append(filesToCompact, *rf)
	}

	if compactFiles, err = s.Compress(dispute.Cid, filesToCompact, path); err != nil {
		return nil, err
	}

	filesInBase64 := base64.StdEncoding.EncodeToString(compactFiles)

	return &protocol.Attachment{
		Name:   fmt.Sprintf("%d.zip", dispute.DisputeID),
		Base64: filesInBase64,
	}, nil
}

func (s svc) Save(chargeback *protocol.Chargeback) error {
	// FIXME: there is a bug here. It must be implement properly
	return s.Save(chargeback)
}
