package chargeback

import (
	"aws-poc/internal/attachment"
	"aws-poc/internal/card"
	"aws-poc/internal/network"
	"aws-poc/internal/protocol"
	"encoding/json"
)

type (
	// A Producer produces protocol.Chargeback through a variety of resources that could be topics, queues and so on.
	Producer interface {
		Produce(*protocol.Chargeback) error
	}

	// A Scheduler schedules protocol.Chargeback status query.
	Scheduler interface {
		Schedule(*protocol.Chargeback) error
	}

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
	// TODO: Refact the organization of this method. Should I declare everything at the top or declare them as soon as I use?
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
		return newIdempotenceError(cid, d.DisputeID)
	}

	if err := s.creator.create(d); err != nil {
		defer s.locker.release(d)
		return newChargebackError(err, cid, d.DisputeID)
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
