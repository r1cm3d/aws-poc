package internal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (
	date              time.Time
	AuthorizationCode string
	ReasonCode        string
	LocalCurrencyCode string
	Queue             string
	Status            string
	Type              string

	Dispute struct {
		Cid       string
		OrgId     string
		AccountId int
		DisputeId int
		AuthorizationCode
		ReasonCode
		CardId            int
		DisputeAmount     float64
		TransactionAmount float64
		TransactionDate   date
		LocalCurrencyCode
		TextMessage       string
		DocumentIndicator bool
		IsPartial         bool
	}

	Chargeback struct {
		*Dispute
		TransactionId string
		ClaimId       string
		ChargebackId  string
		Status
		Queue
		Type
		ResponseError
	}

	locker interface {
		lock(*Dispute) (ok bool)
		release(*Dispute) (ok bool)
	}

	mapper interface {
		fromJSON(string, string) (*Dispute, error)
	}

	creator interface {
		create(*Dispute) error
	}

	opener interface {
		Open(*Dispute, *Card, *Attachment) (*Chargeback, error)
	}

	service struct {
		locker
		mapper
		creator
		cardRegister       CardGetter
		attachmentRegister AttachmentGetter
		opener
	}
)

// ID return DisputeID::CorrelationID
func (e Chargeback) ID() string {
	return fmt.Sprintf("%v::%s", e.DisputeId, e.Cid)
}

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

func (s service) create(dispute *Dispute) error {
	var err error
	var c *Card
	if c, err = s.cardRegister.Get(dispute); err != nil {
		return err
	}
	var att *Attachment
	if att, err = s.attachmentRegister.Get(dispute); err != nil {
		return err
	}

	var cbk *Chargeback
	if cbk, err = s.Open(dispute, c, att); err != nil {
		return err
	}

	fmt.Printf("card: %v", c)
	fmt.Printf("attachment: %v", att)
	fmt.Printf("chargeback: %v", cbk)

	return nil
}

func (s service) handleMessage(cid, body string) error {
	d, err := s.mapper.fromJSON(cid, body)
	if err != nil {
		return newParseError(err)
	}

	if ok := s.locker.lock(d); !ok {
		return newIdempotenceError(cid, d.DisputeId)
	}

	if err := s.creator.create(d); err != nil {
		defer s.locker.release(d)
		return newChargebackError(err, cid, d.DisputeId)
	}

	return nil
}

func (s service) fromJSON(cid, j string) (Dispute, error) {
	var d Dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return Dispute{}, err
	}
	d.Cid = cid
	return d, nil
}
