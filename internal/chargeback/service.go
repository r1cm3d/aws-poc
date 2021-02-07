package chargeback

import (
	"aws-poc/internal/attachment"
	"aws-poc/internal/card"
)

type (
	Type       string
	Queue      string
	ReasonCode string
	Status     string

	Entity struct {
		OrgId         string
		DisputeId     int
		AccountId     int
		TransactionId int
		ClaimId       int
		ChargebackId  int
		Status
		cid string
		Type
		Queue
		Attachment        attachment.Entity
		IsClaimOpen       bool
		DocumentIndicator bool
		RejectReason      string
		ReasonCode
	}

	Input struct {
		Cid               string
		OrgId             string
		DisputeId         int
		AccountId         int
		DocumentIndicator bool
		ReasonCode
		Card       card.Entity
		Attachment attachment.Entity
	}

	Creator interface {
		Create(input Input) (Entity, error)
	}
)
