package protocol

import (
	"fmt"
	"strings"
	"time"
)

type (
	// A File has information about storage files commonly used as attachments.
	File struct {
		Key   string
		Bytes []byte
	}

	// A Date is a wrapper for time.Time.
	Date time.Time

	// An AuthorizationCode is the authorization identifier generated in the authorization process.
	AuthorizationCode string

	// A ReasonCode identifies the type of the chargeback at network brand.
	ReasonCode string

	// A LocalCurrencyCode relates to ISO-4217 standard. It could be Code or Number. For instance: USD or 840.
	// See: https://en.wikipedia.org/wiki/ISO_4217
	LocalCurrencyCode string

	// A Queue relates to network brand queue. This queue is similar to a state.
	Queue string

	// A Status defines if chargeback has error. If it is FAILED has error, otherwise CREATED.
	Status string

	// A Type defines the chargeback type at network brand.
	Type string

	// A NetworkError is used to handle network brand errors.
	NetworkError struct {
	}

	// A Dispute has the input information to create a chargeback at network brand.
	Dispute struct {
		Cid       string
		OrgID     string
		AccountID int
		DisputeID int
		AuthorizationCode
		ReasonCode
		CardID            int
		DisputeAmount     float64
		TransactionAmount float64
		TransactionDate   Date
		LocalCurrencyCode
		TextMessage       string
		DocumentIndicator bool
		IsPartial         bool
	}

	// A Attachment has information about chargeback attachment files.
	Attachment struct {
		Name   string
		Base64 string
	}

	// A Card is used to search for a transaction at network brand.
	// This transaction is used to create a chargeback
	Card struct {
		Number string
	}

	// A Chargeback is the structured that has all information about chargeback.
	Chargeback struct {
		*Dispute
		TransactionID string
		ClaimID       string
		ChargebackID  string
		Status
		Queue
		Type
		*NetworkError
	}
)

// ID returns DisputeID::CorrelationID
func (e Chargeback) ID() string {
	return fmt.Sprintf("%v::%s", e.DisputeID, e.Cid)
}

// ID returns DisputeID::CorrelationID
func (e Dispute) ID() string {
	return fmt.Sprintf("%v::%s", e.DisputeID, e.Cid)
}

// HasError returns true if NetworkError != nil
func (e Chargeback) HasError() bool {
	return e.NetworkError != nil
}

// UnmarshalJSON receives a date in []bytes and parse it in the pattern YYYY-MM-DD
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*d = Date(t)

	return nil
}

func (n NetworkError) Error() string {
	return fmt.Sprintf("network error: ")
}
