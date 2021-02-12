package chargeback

import (
	"aws-poc/internal/protocol"
	"errors"
)
//TODO: move it to protocol package
const (
	disputeID         = 666
	cid               = "e1388e36-1683-4902-b30c-5c5b63f5976c"
	orgId             = "TN-ed3d9cbf-664e-4044-bc1f-5adee7ff069f"
	accountId         = 10782
	transactionId     = "26811"
	claimId           = "5717"
	chargebackId      = "27202"
	cardId            = 27542
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
	stubError     = errors.New("mocked error")
	cardError     = errors.New("mocked card error")
	attGetError   = errors.New("mocked att get error")
	openerError   = errors.New("mocked opener error")
	producerError = errors.New("mocked producer error")
	scdError      = errors.New("mocked scheduler error")
	attSaveError  = errors.New("mocked att save error")
	disputeStub   = &protocol.Dispute{
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

	chargebackStub = &protocol.Chargeback{
		Dispute:       disputeStub,
		TransactionId: transactionId,
		ClaimId:       claimId,
		ChargebackId:  chargebackId,
		Status:        protocol.Status("CREATED"),
		Queue:         protocol.Queue("REJECTS"),
		Type:          protocol.Type("CHARGEBACK"),
		NetworkError:  nil,
	}

	chargebackWithErrorStub = &protocol.Chargeback{
		Dispute:       disputeStub,
		TransactionId: transactionId,
		ClaimId:       claimId,
		ChargebackId:  chargebackId,
		Status:        protocol.Status("CREATED"),
		Queue:         protocol.Queue("REJECTS"),
		Type:          protocol.Type("CHARGEBACK"),
		NetworkError:  &protocol.NetworkError{},
	}

	cardStub       = &protocol.Card{Number: "5172163143182969"}
	attachmentStub = &protocol.Attachment{Name: "filename", Base64: "ZmlsZW5hbWUgaW4gYmFzZTY0"}
)
