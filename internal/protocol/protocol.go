package protocol

import (
	"fmt"
	"strings"
	"time"
)

type (
	File struct {
		Key string
	}

	Date              time.Time
	AuthorizationCode string
	ReasonCode        string
	LocalCurrencyCode string
	Queue             string
	Status            string
	Type              string

	NetworkError struct {
	}

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

	Attachment struct {
		Name   string
		Base64 string
	}

	Card struct {
		Number string
	}

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

// ID return DisputeID::CorrelationID
func (e Chargeback) ID() string {
	return fmt.Sprintf("%v::%s", e.DisputeID, e.Cid)
}

// ID return DisputeID::CorrelationID
func (e Dispute) ID() string {
	return fmt.Sprintf("%v::%s", e.DisputeID, e.Cid)
}

// HasError return true if NetworkError != nil
func (e Chargeback) HasError() bool {
	return e.NetworkError != nil
}

// UnmarshalJSON receive a date in []bytes and parse it in the pattern YYYY-MM-DD
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
