package attachment

import (
	"aws-poc/internal/protocol"
	"errors"
	"fmt"
)

const (
	disputeID         = 777
	cid               = "6489f479-b609-4b3d-ab15-d5947f012c3c"
	orgId             = "TN-7b90a75d-d094-4498-a3c8-7cec480f216f"
	accountId         = 10787
	cardId            = 28542
	disputeAmount     = 150.00
	transactionAmount = 200.00
	documentIndicator = false
	reasonCode        = protocol.ReasonCode("4853")
	authorizationCode = protocol.AuthorizationCode("AB6ZZR")
	usDollar          = protocol.LocalCurrencyCode("840")
	isPartial         = false
	textMessage       = "Example message"
	transactionId     = "26811"
	claimId           = "5717"
	chargebackId      = "27202"
)

var (
	listError        = errors.New("storage list error")
	getError         = errors.New("storage get error")
	unsentFilesError = errors.New("unsent files error")
	archiverError    = errors.New("archiver error")
	saveError        = errors.New("save error")
	path             = fmt.Sprintf("%s/%d/%d", filenameRoot, disputeStub.AccountId, disputeStub.DisputeId)
	f1               = file{key: "cbk_file1.pdf"}
	f2               = file{key: "cbk_doc.pdf"}
	f3               = file{key: "file3.pdf"}
	fg1              = file{key: "cbk_get_file1.pdf"}
	fg2              = file{key: "cbk_get_doc.pdf"}
	fg3              = file{key: "file_get_3.pdf"}
	uf1              = file{key: fmt.Sprintf("%s/%s", path, f1.key)}
	uf2              = file{key: fmt.Sprintf("%s/%s", path, f2.key)}
	uf3              = file{key: fmt.Sprintf("%s/%s", path, f3.key)}
	files            = []file{f1, f2, f3}
	unsentFiles      = []file{uf1, uf2, uf3}
	getFiles         = []file{fg1, fg2, fg3}
	attStub          = &protocol.Attachment{Name: "cbk666.zip", Base64: "ZmlsZW5hbWUgaW4gYmFzZTY0"}
	disputeStub      = &protocol.Dispute{
		Cid:               cid,
		OrgId:             orgId,
		AccountId:         accountId,
		DisputeId:         disputeID,
		AuthorizationCode: authorizationCode,
		ReasonCode:        reasonCode,
		CardId:            cardId,
		DisputeAmount:     disputeAmount,
		TransactionAmount: transactionAmount,
		LocalCurrencyCode: usDollar,
		TextMessage:       textMessage,
		DocumentIndicator: documentIndicator,
		IsPartial:         isPartial,
	}
	storageStub = &mockStorage{
		expPath:  path,
		expFiles: [3][2]file{{uf1, fg1}, {uf2, fg2}, {uf3, fg3}},
	}
	chargebackStub = &protocol.Chargeback{
		Dispute:       disputeStub,
		TransactionId: transactionId,
		ClaimId:       claimId,
		ChargebackId:  chargebackId,
		Status:        protocol.Status("CREATED"),
		Queue:         protocol.Queue("CLOSED"),
		Type:          protocol.Type("SECOND_PRESENTMENT"),
		NetworkError:  nil,
	}
)
