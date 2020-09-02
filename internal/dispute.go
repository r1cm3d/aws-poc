package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type (
	Date time.Time

	Dispute struct {
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
	disputeRepository interface {
		lock(string, Dispute) (ok bool)
		unlock(string, Dispute)
	}

	disputeMapper interface {
		mapFromJson(string) (Dispute, error)
	}

	disputeSvc struct {
		disputeRepository
		disputeMapper
	}
)

func (s disputeSvc) openChargeback(_ Dispute) error {
	return nil
}

func (s disputeSvc) handleMessage(cid, body string) error {
	d, err := s.mapFromJson(body) // TODO: extract it to an interface to be able to test it
	if err != nil {
		return errors.New(fmt.Sprintf("parser error: %s", err.Error()))
	}

	if ok := s.disputeRepository.lock(cid, d); !ok { //TODO: cover this flow
		return errors.New(fmt.Sprintf("idempotence error: cid(%v), disputeId(%v)", cid, d.DisputeId))
	}

	if err := s.openChargeback(d); err != nil { //TODO: cover this flow
		defer s.disputeRepository.unlock(cid, d)
		return errors.New(fmt.Sprintf("parser error: %s", err.Error()))
	}

	return nil
}

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

func (s disputeSvc) mapFromJson(j string) (Dispute, error) {
	var d Dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return Dispute{}, err
	}
	return d, nil
}
