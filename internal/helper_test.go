package internal

import (
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
	reasonCode        = ReasonCode("4853")
	authorizationCode = AuthorizationCode("ABDZAR")
	usDollar          = LocalCurrencyCode("840")
	isPartial         = false
	textMessage       = "Example message"
)

var (
	errStub     = errors.New("mocked error")
	disputeStub = Dispute{
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

	chargebackStub = Chargeback{
		Dispute:       disputeStub,
		TransactionId: transactionId,
		ClaimId:       claimId,
		ChargebackId:  chargebackId,
		Status:        Status("CREATED"),
		Queue:         Queue("REJECTS"),
		Type:          Type("CHARGEBACK"),
		ResponseError: ResponseError{},
	}

	cardStub       = Card{Number: "5172163143182969"}
	attachmentStub = Attachment{Name: "filename", Base64: "ZmlsZW5hbWUgaW4gYmFzZTY0"}
)
