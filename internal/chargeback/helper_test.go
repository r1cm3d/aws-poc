package chargeback

import (
	"aws-poc/internal/protocol"
	"errors"
)

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
	stubError   = errors.New("mocked error")
	cardError   = errors.New("mocked card error")
	attError    = errors.New("mocked att error")
	openerError = errors.New("mocked opener error")
	disputeStub = &protocol.Dispute{
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
		ResponseError: protocol.ResponseError{},
	}

	cardStub       = &protocol.Card{Number: "5172163143182969"}
	attachmentStub = &protocol.Attachment{Name: "filename", Base64: "ZmlsZW5hbWUgaW4gYmFzZTY0"}
)
