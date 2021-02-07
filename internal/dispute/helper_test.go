package dispute

import (
	"aws-poc/internal/attachment"
	"aws-poc/internal/card"
	"aws-poc/internal/chargeback"
	"errors"
)

const (
	disputeID         = 666
	cid               = "e1388e36-1683-4902-b30c-5c5b63f5976c"
	orgId             = "TN-ed3d9cbf-664e-4044-bc1f-5adee7ff069f"
	accountId         = 10782
	transactionId     = 26811
	claimId           = 5717
	chargebackId      = 27202
	cardId            = 027542
	disputeAmount     = 120.00
	documentIndicator = false
	reasonCode        = ReasonCode("4853")
)

var (
	errFake     = errors.New("mocked error")
	disputeFake = Entity{
		CorrelationID:       cid,
		DisputeID:           disputeID,
		AccountID:           accountId,
		AuthorizationCode:   AuthorizationCode("JT11F6"),
		ReasonCode:          reasonCode,
		CardID:              cardId,
		Tenant:              orgId,
		DisputeAmount:       disputeAmount,
		TransactionDate:     date{},
		LocalCurrencyCode:   LocalCurrencyCode("986"),
		TextMessage:         "Chargeback test",
		DocumentIndicator:   documentIndicator,
		IsPartialChargeback: false,
	}
	cardFake       = card.Entity{Number: "5172163143182969"}
	attachmentFake = attachment.Entity{Name: "filename", Base64: "ZmlsZW5hbWUgaW4gYmFzZTY0"}
	chargebackFake = chargeback.Entity{
		OrgId:             orgId,
		DisputeId:         disputeID,
		AccountId:         accountId,
		TransactionId:     transactionId,
		ClaimId:           claimId,
		ChargebackId:      chargebackId,
		Status:            chargeback.Status("ERROR"),
		Type:              chargeback.Type("CHARGEBACK"),
		Queue:             chargeback.Queue("REJECTS"),
		Attachment:        attachment.Entity{},
		IsClaimOpen:       false,
		DocumentIndicator: false,
		RejectReason:      "Code1=0142(00):D0063/002;DE072=D0063\\\\\\\\8000000808\\\\\\\\\\\\\\\\",
		ReasonCode:        chargeback.ReasonCode("4853"),
	}
)
