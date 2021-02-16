package awsrepository

import (
	"aws-poc/internal/protocol"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
)

const (
	disputeID         = 666
	cid               = "e1388e36-1683-4902-b30c-5c5b63f5976c"
	orgId             = "TN-ed3d9cbf-664e-4044-bc1f-5adee7ff069f"
	accountId         = 10782
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
	tableName                   = aws.String("Dispute")
	id                          = aws.String("ID")
	parserError                 = errors.New("parseError")
	putItemError                = errors.New("putItemError")
	deleteError                 = errors.New("deleteError")
	queryError                  = errors.New("queryError")
	getError					= errors.New("getError")
	unmarshallerListOfMapsError = errors.New("UnmarshallerListOfMapsError")
	unmarshallerError           = errors.New("unmarshallerError")

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
)

func defaultInput() record {
	return Item{
		DisputeID: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}
}