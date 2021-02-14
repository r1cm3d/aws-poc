package chargeback

import (
	"aws-poc/internal/attachment"
	"aws-poc/internal/card"
	"aws-poc/internal/network"
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

	Producer interface {
		Produce(*protocol.Chargeback) error
	}

	Scheduler interface {
		Schedule(*protocol.Chargeback) error
	}

	svc struct {
		locker
		mapper
		creator
		cardService    card.Service
		attService     attachment.Service
		networkCreator network.Creator
		Scheduler
		Producer
	}
)

func (s svc) create(dispute *protocol.Dispute) error {
	var err error
	var c *protocol.Card
	if c, err = s.cardService.Get(dispute); err != nil {
		return err
	}
	var att *protocol.Attachment
	if att, err = s.attService.Get(dispute); err != nil {
		return err
	}
	var cbk *protocol.Chargeback
	if cbk, err = s.networkCreator.Create(dispute, c, att); err != nil {
		return err
	}
	if err = s.Produce(cbk); err != nil {
		return err
	}
	if cbk.HasError() {
		return cbk.NetworkError
	}
	if err = s.attService.Save(cbk); err != nil {
		return err
	}
	if err = s.Schedule(cbk); err != nil {
		return err
	}

	return nil
}

func (s svc) handleMessage(cid, body string) error {
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

func (s svc) fromJSON(cid, j string) (protocol.Dispute, error) {
	var d protocol.Dispute
	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		return protocol.Dispute{}, err
	}
	d.Cid = cid
	return d, nil
}
