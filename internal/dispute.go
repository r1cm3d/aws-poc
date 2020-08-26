package internal

import (
	"encoding/json"
	"strings"
	"time"
)

type Dispute struct {
	CorrelationId       string
	DisputeId           int
	AccountId           int
	AuthorizationCode   string
	ReasonCode          string
	CardId              string
	Tenant              string
	DisputeAmount       float64
	TransactionDate     Date
	LocalCurrencyCode   string
	TextMessage         string
	DocumentIndicator   bool
	IsPartialChargeback bool
}

type Date time.Time

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

func mapFromJson(j string) (Dispute, error) {
	var d Dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return Dispute{}, err
	}
	return d, nil
}
