package chargeback

import (
	"aws-poc/internal/attachment"
	"aws-poc/internal/card"
	"aws-poc/internal/protocol"
	"encoding/json"
)

type (
	locker interface {
		lock(*protocol.Dispute) (ok bool)
		release(*protocol.Dispute) (ok bool)
	}

	mapper interface {
		fromJSON(string, string) (*protocol.Dispute, error)
	}

	creator interface {
		create(*protocol.Dispute) error
	}

	opener interface {
		Open(*protocol.Dispute, *protocol.Card, *protocol.Attachment) (*protocol.Chargeback, error)
	}

	Producer interface {
		Produce(*protocol.Chargeback) error
	}

	Scheduler interface {
		Schedule(*protocol.Chargeback) error
	}

	service struct {
		locker
		mapper
		creator
		cardRepository       card.Repository
		attachmentRepository attachment.Repository
		opener
		Scheduler
		Producer
	}
)

func (s service) create(dispute *protocol.Dispute) error {
	var err error
	var c *protocol.Card
	if c, err = s.cardRepository.Get(dispute); err != nil {
		return err
	}
	var att *protocol.Attachment
	if att, err = s.attachmentRepository.Get(dispute); err != nil {
		return err
	}
	var cbk *protocol.Chargeback
	if cbk, err = s.Open(dispute, c, att); err != nil {
		return err
	}
	if err = s.Produce(cbk); err != nil {
		return err
	}
	if cbk.HasError() {
		return cbk.NetworkError
	}
	if err = s.attachmentRepository.Save(cbk); err != nil {
		return err
	}
	if err = s.Schedule(cbk); err != nil {
		return err
	}

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

func (s service) fromJSON(cid, j string) (protocol.Dispute, error) {
	var d protocol.Dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return protocol.Dispute{}, err
	}
	d.Cid = cid
	return d, nil
}
