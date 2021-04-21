package chargeback

import (
	"aws-poc/internal/protocol"
	"reflect"
	"testing"
	"time"
)

func TestMapFromJson(t *testing.T) {
	svc := svc{}
	cid := cid
	json := `{
  "disputeId": 611,
  "accountID": 48448,
  "authorizationCode": "7HSGXW",
  "reasonCode": "848",
  "cardID": 3123,
  "orgID": "pismo.io",
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
		DisputeID:         611,
		AccountID:         48448,
		AuthorizationCode: protocol.AuthorizationCode("7HSGXW"),
		ReasonCode:        protocol.ReasonCode("848"),
		CardID:            3123,
		OrgID:             "pismo.io",
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
	svc := svc{}

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
		svc
	}{
		{"success", defaultInput, nil, svc{mapper: mockMapper{}, locker: mockRepository{}, creator: mockCreator{}}},
		{"parseError", defaultInput, newParseError(errStub), svc{locker: mockRepository{}, mapper: errMapper{}}},
		{"idempotenceError", defaultInput, newIdempotenceError(cid, disputeID), svc{mapper: mockMapper{}, locker: errRepository{}}},
		{"chargebackError", defaultInput, newChargebackError(errStub, cid, disputeID), svc{mapper: mockMapper{}, locker: mockRepository{}, creator: errDisputer{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.svc.handleMessage(c.in[0], c.in[1]); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestCreateSuccess(t *testing.T) {
	cr, ar, ope, prod, scd := mockCardService{}, mockAttService{}, mockNetworkCreator{}, mockProducer{expected: chargebackStub}, mockScheduler{}
	svc := svc{
		cardService:    &cr,
		attService:     &ar,
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
	cr, ar, ope, prod, scd := mockCardService{}, mockAttService{}, mockOpenerWithNetworkError{}, mockProducer{expected: chargebackWithErrorStub}, mockScheduler{}
	svc := svc{
		cardService:    &cr,
		attService:     &ar,
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
		svc  svc
		want error
	}{
		{"errCardStub", disputeStub, svc{
			cardService: errCardService{},
		}, errCardStub},
		{"attachmentGetError", disputeStub, svc{
			cardService: &mockCardService{},
			attService:  &errAttGetService{},
		}, errAttGetStub},
		{"errOpenerStub", disputeStub, svc{
			cardService:    &mockCardService{},
			attService:     &mockAttService{},
			networkCreator: errNetworkCreator{},
		}, errOpenerStub},
		{"errProducerStub", disputeStub, svc{
			cardService:    &mockCardService{},
			attService:     &mockAttService{},
			networkCreator: &mockNetworkCreator{},
			Producer:       &errProducer{},
		}, errProducerStub},
		{"attachmentSaveError", disputeStub, svc{
			cardService:    &mockCardService{},
			attService:     &errAttSaveService{},
			networkCreator: &mockNetworkCreator{},
			Producer:       &mockProducer{},
		}, errAttSaveStub},
		{"scheduleError", disputeStub, svc{
			cardService:    &mockCardService{},
			attService:     &mockAttService{},
			networkCreator: &mockNetworkCreator{},
			Producer:       &mockProducer{},
			Scheduler:      &errScheduler{},
		}, errScdStub},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.svc.create(c.in); got != c.want {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
