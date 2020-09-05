package internal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (
	// Date is an type to parse a string in the format YYYY-MM-DD and parse to time.Time
	Date time.Time

	dispute struct {
		CorrelationID       string
		DisputeID           int
		AccountID           int
		AuthorizationCode   string
		ReasonCode          string
		CardID              string
		Tenant              string
		DisputeAmount       float64
		TransactionDate     Date
		LocalCurrencyCode   string
		TextMessage         string
		DocumentIndicator   bool
		IsPartialChargeback bool
	}
	disputeRepository interface {
		lock(dispute) (ok bool)
		unlock(dispute)
	}

	disputeMapper interface {
		mapFromJSON(string, string) (dispute, error)
	}

	disputer interface {
		open(dispute) error
	}

	disputeSvc struct {
		disputeRepository
		disputeMapper
		disputer
	}
)

// UnmarshalJSON receive a date in []bytes and parse it in the pattern YYYY-MM-DD
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" { //TODO: cover this flow
		return nil
	}

	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil { //TODO: cover this flow
		return err
	}

	*d = Date(t)

	return nil
}

func (s disputeSvc) open(_ dispute) error {
	return nil
}

func (s disputeSvc) handleMessage(cid, body string) error {
	d, err := s.mapFromJSON(cid, body)
	if err != nil {
		return fmt.Errorf("parser error: %s", err.Error())
	}

	if ok := s.disputeRepository.lock(d); !ok {
		return fmt.Errorf("idempotence error: cid(%v), disputeId(%v)", cid, d.DisputeID)
	}

	if err := s.open(d); err != nil {
		defer s.disputeRepository.unlock(d)
		return fmt.Errorf("parser error: %s", err.Error())
	}

	return nil
}

func (s disputeSvc) mapFromJSON(cid, j string) (dispute, error) {
	var d dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return dispute{}, err
	}
	d.CorrelationID = cid
	return d, nil
}
