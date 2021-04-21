package chargeback

import (
	"aws-poc/internal/protocol"
	"errors"
)

const (
	disputeID         = 666
	cid               = "e1388e36-1683-4902-b30c-5c5b63f5976c"
	orgID             = "TN-ed3d9cbf-664e-4044-bc1f-5adee7ff069f"
	accountID         = 10782
	transactionID     = "26811"
	claimID           = "5717"
	chargebackID      = "27202"
	cardID            = 27542
	disputeAmount     = 120.00
	transactionAmount = 150.00
	documentIndicator = false
	reasonCode        = protocol.ReasonCode("4853")
	authorizationCode = protocol.AuthorizationCode("ABDZAR")
	usDollar          = protocol.LocalCurrencyCode("840")
	isPartial         = false
	textMessage       = "Example message"
)

var (
	errStub         = errors.New("mocked error")
	errCardStub     = errors.New("mocked card error")
	errAttGetStub   = errors.New("mocked att get error")
	errOpenerStub   = errors.New("mocked opener error")
	errProducerStub = errors.New("mocked producer error")
	errScdStub      = errors.New("mocked scheduler error")
	errAttSaveStub  = errors.New("mocked att save error")
	disputeStub     = &protocol.Dispute{
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

	chargebackStub = &protocol.Chargeback{
		Dispute:       disputeStub,
		TransactionID: transactionID,
		ClaimID:       claimID,
		ChargebackID:  chargebackID,
		Status:        protocol.Status("CREATED"),
		Queue:         protocol.Queue("REJECTS"),
		Type:          protocol.Type("CHARGEBACK"),
		NetworkError:  nil,
	}

	chargebackWithErrorStub = &protocol.Chargeback{
		Dispute:       disputeStub,
		TransactionID: transactionID,
		ClaimID:       claimID,
		ChargebackID:  chargebackID,
		Status:        protocol.Status("CREATED"),
		Queue:         protocol.Queue("REJECTS"),
		Type:          protocol.Type("CHARGEBACK"),
		NetworkError:  &protocol.NetworkError{},
	}

	cardStub       = &protocol.Card{Number: "5172163143182969"}
	attachmentStub = &protocol.Attachment{Name: "filename", Base64: "ZmlsZW5hbWUgaW4gYmFzZTY0"}
)
