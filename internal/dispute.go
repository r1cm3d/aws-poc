package internal

import (
	"encoding/json"
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
	LocalCurrencyCode   string
	TextMessage         string
	DocumentIndicator   bool
	IsPartialChargeback bool
}

func mapFromJson(j string) (Dispute, error) {
	var d Dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return Dispute{}, err
	}
	return d, nil
}
