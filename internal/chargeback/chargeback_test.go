package chargeback

import (
	"aws-poc/internal/protocol"
	"reflect"
	"testing"
	"time"
)

func TestMapFromJson(t *testing.T) {
	svc := service{}
	cid := cid
	json := `{
  "disputeId": 611,
  "accountId": 48448,
  "authorizationCode": "7HSGXW",
  "reasonCode": "848",
  "cardId": 3123,
  "orgId": "pismo.io",
  "disputeAmount": 32.32,
  "transactionAmount": 42.65,
  "transactionDate": "2012-04-23",
  "localCurrencyCode": "986",
  "textMessage": "this a test message",
  "documentIndicator": true,
  "isPartial": false
}`
	want := protocol.Dispute{
		Cid:               cid,
		DisputeId:         611,
		AccountId:         48448,
		AuthorizationCode: protocol.AuthorizationCode("7HSGXW"),
		ReasonCode:        protocol.ReasonCode("848"),
		CardId:            3123,
		OrgId:             "pismo.io",
		DisputeAmount:     32.32,
		TransactionDate:   protocol.Date(time.Date(2012, 04, 23, 0, 0, 0, 0, time.UTC)),
		LocalCurrencyCode: protocol.LocalCurrencyCode("986"),
		TextMessage:       "this is a test message",
		DocumentIndicator: true,
		IsPartial:         false,
	}

	got, err := svc.fromJSON(cid, json)

	if err != nil {
		t.Error("fromJSON() error should not be returned")
	}
	if reflect.DeepEqual(got, want) {
		t.Errorf("fromJSON() got: %v, want: %v", got, want)
	}
}

func TestMapFromJson_Error(t *testing.T) {
	svc := service{}

	_, err := svc.fromJSON("", "json")

	if err == nil {
		t.Error("mapFromJson_Error() error should be returned")
	}
}

func TestHandleMessage(t *testing.T) {
	defaultInput := [2]string{cid, "body"}
	cases := []struct {
		name string
		in   [2]string
		want error
		service
	}{
		{"success", defaultInput, nil, service{mapper: mockMapper{}, locker: mockRepository{}, creator: mockCreator{}}},
		{"parseError", defaultInput, newParseError(stubError), service{locker: mockRepository{}, mapper: errMapper{}}},
		{"idempotenceError", defaultInput, newIdempotenceError(cid, disputeID), service{mapper: mockMapper{}, locker: errRepository{}}},
		{"chargebackError", defaultInput, newChargebackError(stubError, cid, disputeID), service{mapper: mockMapper{}, locker: mockRepository{}, creator: errDisputer{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.service.handleMessage(c.in[0], c.in[1]); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestCreateSuccess(t *testing.T) {
	cr, ar, ope, prod, scd := mockCardRepository{}, mockAttachmentRepository{}, mockNetworkCreator{}, mockProducer{expected: chargebackStub}, mockScheduler{}
	svc := service{
		cardRepository: &cr,
		attRepository:  &ar,
		networkCreator: &ope,
		Producer:       &prod,
		Scheduler:      &scd,
	}
	_ = svc.create(disputeStub)

	if !cr.called {
		t.Error("card register not called")
	}
	if !ar.getCalled {
		t.Error("attachment register get not called")
	}
	if !ope.called {
		t.Error("chargeback creator not called")
	}
	if !prod.called {
		t.Error("chargeback producer not called")
	}
	if !ar.saveCalled {
		t.Error("attachment register save not called")
	}
	if !scd.called {
		t.Errorf("scheduler not called")
	}
}

func TestCreateNetworkError(t *testing.T) {
	cr, ar, ope, prod, scd := mockCardRepository{}, mockAttachmentRepository{}, mockOpenerWithNetworkError{}, mockProducer{expected: chargebackWithErrorStub}, mockScheduler{}
	svc := service{
		cardRepository: &cr,
		attRepository:  &ar,
		networkCreator: &ope,
		Producer:       &prod,
		Scheduler:      &scd,
	}
	_ = svc.create(disputeStub)

	if !cr.called {
		t.Error("card register not called")
	}
	if !ar.getCalled {
		t.Error("attachment register get not called")
	}
	if !ope.called {
		t.Error("chargeback creator not called")
	}
	if !prod.called {
		t.Error("chargeback producer not called")
	}
	if ar.saveCalled {
		t.Error("attachment register save called")
	}
	if scd.called {
		t.Errorf("scheduler called")
	}
}

func TestOpenFail(t *testing.T) {
	cases := []struct {
		name string
		in   *protocol.Dispute
		svc  service
		want error
	}{
		{"cardError", disputeStub, service{
			cardRepository: errCardGetter{},
		}, cardError},
		{"attachmentGetError", disputeStub, service{
			cardRepository: &mockCardRepository{},
			attRepository:  &errAttachmentGetRepository{},
		}, attGetError},
		{"openerError", disputeStub, service{
			cardRepository: &mockCardRepository{},
			attRepository:  &mockAttachmentRepository{},
			networkCreator: errNetworkCreator{},
		}, openerError},
		{"producerError", disputeStub, service{
			cardRepository: &mockCardRepository{},
			attRepository:  &mockAttachmentRepository{},
			networkCreator: &mockNetworkCreator{},
			Producer:       &errProducer{},
		}, producerError},
		{"attachmentSaveError", disputeStub, service{
			cardRepository: &mockCardRepository{},
			attRepository:  &errAttachmentSaveRepository{},
			networkCreator: &mockNetworkCreator{},
			Producer:       &mockProducer{},
		}, attSaveError},
		{"scheduleError", disputeStub, service{
			cardRepository: &mockCardRepository{},
			attRepository:  &mockAttachmentRepository{},
			networkCreator: &mockNetworkCreator{},
			Producer:       &mockProducer{},
			Scheduler:      &errScheduler{},
		}, scdError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.svc.create(c.in); got != c.want {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
