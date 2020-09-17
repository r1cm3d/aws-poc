package dispute

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	date time.Time

	dispute struct {
		CorrelationID       string
		DisputeID           int
		AccountID           int
		AuthorizationCode   string
		ReasonCode          string
		CardID              string
		Tenant              string
		DisputeAmount       float64
		TransactionDate     date
		LocalCurrencyCode   string
		TextMessage         string
		DocumentIndicator   bool
		IsPartialChargeback bool
	}

	register interface {
		lock(dispute) (ok bool)
		unlock(dispute)
	}

	mapper interface {
		fromJSON(string, string) (dispute, error)
	}

	disputer interface {
		open(dispute) error
	}

	service struct {
		register
		mapper
		disputer
	}
)

// UnmarshalJSON receive a date in []bytes and parse it in the pattern YYYY-MM-DD
func (d *date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*d = date(t)

	return nil
}

func (s service) open(_ dispute) error {
	return nil
}

func (s service) handleMessage(cid, body string) error {
	d, err := s.mapper.fromJSON(cid, body)
	if err != nil {
		return newParseError(err)
	}

	if ok := s.register.lock(d); !ok {
		return newIdempotenceError(cid, d.DisputeID)
	}

	if err := s.disputer.open(d); err != nil {
		defer s.register.unlock(d)
		return newChargebackError(err, cid, d.DisputeID)
	}

	return nil
}

func (s service) fromJSON(cid, j string) (dispute, error) {
	var d dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return dispute{}, err
	}
	d.CorrelationID = cid
	return d, nil
}
