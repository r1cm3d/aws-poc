package attachment

import (
	"aws-poc/internal/protocol"
	"errors"
	"fmt"
)

const (
	disputeID         = 777
	cid               = "6489f479-b609-4b3d-ab15-d5947f012c3c"
	orgID             = "TN-7b90a75d-d094-4498-a3c8-7cec480f216f"
	accountID         = 10787
	cardID            = 28542
	disputeAmount     = 150.00
	transactionAmount = 200.00
	documentIndicator = false
	reasonCode        = protocol.ReasonCode("4853")
	authorizationCode = protocol.AuthorizationCode("AB6ZZR")
	usDollar          = protocol.LocalCurrencyCode("840")
	isPartial         = false
	textMessage       = "Example message"
	transactionID     = "26811"
	claimID           = "5717"
	chargebackID      = "27202"
)

var (
	errListStub        = errors.New("storage list error")
	errGetStub         = errors.New("storage Get error")
	errUnsentFilesStub = errors.New("unsent files error")
	errArchiverStub    = errors.New("Compressor error")
	errSaveStub        = errors.New("save error")
	path               = fmt.Sprintf("%s/%d/%d", filenameRoot, disputeStub.AccountID, disputeStub.DisputeID)
	f1                 = protocol.File{Key: "cbk_file1.pdf"}
	f2                 = protocol.File{Key: "cbk_doc.pdf"}
	f3                 = protocol.File{Key: "file3.pdf"}
	fg1                = protocol.File{Key: "cbk_get_file1.pdf"}
	fg2                = protocol.File{Key: "cbk_get_doc.pdf"}
	fg3                = protocol.File{Key: "file_get_3.pdf"}
	uf1                = protocol.File{Key: fmt.Sprintf("%s/%s", path, f1.Key)}
	uf2                = protocol.File{Key: fmt.Sprintf("%s/%s", path, f2.Key)}
	uf3                = protocol.File{Key: fmt.Sprintf("%s/%s", path, f3.Key)}
	files              = []protocol.File{f1, f2, f3}
	unsentFiles        = []protocol.File{uf1, uf2, uf3}
	getFiles           = []protocol.File{fg1, fg2, fg3}
	attStub            = &protocol.Attachment{Name: "777.zip", Base64: "Wm1sc1pXNWhiV1VnYVc0Z1ltRnpaVFkw"}
	compactFilesStub   = []byte("ZmlsZW5hbWUgaW4gYmFzZTY0")
	disputeStub        = &protocol.Dispute{
		Cid:               cid,
		OrgID:             orgID,
		AccountID:         accountID,
		DisputeID:         disputeID,
		AuthorizationCode: authorizationCode,
		ReasonCode:        reasonCode,
		CardID:            cardID,
		DisputeAmount:     disputeAmount,
		TransactionAmount: transactionAmount,
		LocalCurrencyCode: usDollar,
		TextMessage:       textMessage,
		DocumentIndicator: documentIndicator,
		IsPartial:         isPartial,
	}
	storageStub = &mockStorage{
		expPath:  path,
		expFiles: [3][2]protocol.File{{uf1, fg1}, {uf2, fg2}, {uf3, fg3}},
	}
	chargebackStub = &protocol.Chargeback{
		Dispute:       disputeStub,
		TransactionID: transactionID,
		ClaimID:       claimID,
		ChargebackID:  chargebackID,
		Status:        protocol.Status("CREATED"),
		Queue:         protocol.Queue("CLOSED"),
		Type:          protocol.Type("SECOND_PRESENTMENT"),
		NetworkError:  nil,
	}
)
