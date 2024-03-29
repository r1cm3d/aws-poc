package awsrepository

import (
	"aws-poc/internal/protocol"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
)

const (
	disputeID         = 666
	cid               = "e1388e36-1683-4902-b30c-5c5b63f5976c"
	orgID             = "TN-ed3d9cbf-664e-4044-bc1f-5adee7ff069f"
	accountID         = 10782
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
	tableName                     = aws.String("Dispute")
	id                            = aws.String("ID")
	errParserStub                 = errors.New("parseError")
	errPutItemStub                = errors.New("errPutItemStub")
	errDeleteStub                 = errors.New("errDeleteStub")
	errQueryStub                  = errors.New("errQueryStub")
	errGetStub                    = errors.New("errGetStub")
	errUnmarshallerListOfMapsStub = errors.New("UnmarshallerListOfMapsError")
	errUnmarshallerStub           = errors.New("errUnmarshallerStub")

	disputeStub = &protocol.Dispute{
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
)

func defaultInput() record {
	return Item{
		DisputeID: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}
}
