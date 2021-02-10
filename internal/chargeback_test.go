package internal

import (
	"reflect"
	"testing"
	"time"
)

type (
	errMapper            struct{}
	errRepository        struct{}
	errDisputer          struct{}
	mockRepository       struct{}
	mockMapper           struct{}
	mockCreator          struct{}
	mockCardGetter       struct{}
	mockAttachmentGetter struct{}
	mockOpener           struct{}
	errCardGetter       struct{}
	errAttachmentGetter struct{}
	errOpener           struct{}
)

var (
	cardRegisterCalled       bool
	attachmentRegisterCalled bool
	chargebackCreatorCalled  bool
)

func (e errMapper) fromJSON(string, string) (*Dispute, error) {
	return disputeStub, errStub
}

func (m mockMapper) fromJSON(string, string) (*Dispute, error) {
	return disputeStub, nil
}

func (e errRepository) lock(*Dispute) (ok bool) {
	return false
}

func (e errRepository) release(*Dispute) (ok bool) {
	return false
}

func (e errDisputer) create(*Dispute) error {
	return errStub
}

func (m mockCreator) create(*Dispute) error {
	return nil
}

func (m mockRepository) lock(*Dispute) (ok bool) {
	return true
}

func (m mockRepository) release(*Dispute) (ok bool) {
	return true
}

func (m mockCardGetter) Get(dispute *Dispute) (*Card, error) {
	cardRegisterCalled = dispute == disputeStub

	return cardStub, nil
}

func (m mockAttachmentGetter) Get(dispute *Dispute) (*Attachment, error) {
	attachmentRegisterCalled = dispute == disputeStub

	return attachmentStub, nil
}

func (m mockOpener) Open(dispute *Dispute, card *Card, attachment *Attachment) (*Chargeback, error) {
	chargebackCreatorCalled = dispute == disputeStub &&
		card == cardStub &&
		attachment == attachmentStub

	return chargebackStub, nil
}

func (m errCardGetter) Get(dispute *Dispute) (*Card, error) {
	cardRegisterCalled = dispute == disputeStub

	return nil, errStub
}

func (m errAttachmentGetter) Get(dispute *Dispute) (*Attachment, error) {
	attachmentRegisterCalled = dispute == disputeStub

	return attachmentStub, nil
}

func (m errOpener) Open(dispute *Dispute, card *Card, attachment *Attachment) (*Chargeback, error) {
	chargebackCreatorCalled = dispute == disputeStub &&
		card == cardStub &&
		attachment == attachmentStub

	return chargebackStub, nil
}

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
	want := Dispute{
		Cid:               cid,
		DisputeId:         611,
		AccountId:         48448,
		AuthorizationCode: AuthorizationCode("7HSGXW"),
		ReasonCode:        ReasonCode("848"),
		CardId:            3123,
		OrgId:             "pismo.io",
		DisputeAmount:     32.32,
		TransactionDate:   date(time.Date(2012, 04, 23, 0, 0, 0, 0, time.UTC)),
		LocalCurrencyCode: LocalCurrencyCode("986"),
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
		{"parseError", defaultInput, newParseError(errStub), service{locker: mockRepository{}, mapper: errMapper{}}},
		{"idempotenceError", defaultInput, newIdempotenceError(cid, disputeID), service{mapper: mockMapper{}, locker: errRepository{}}},
		{"chargebackError", defaultInput, newChargebackError(errStub, cid, disputeID), service{mapper: mockMapper{}, locker: mockRepository{}, creator: errDisputer{}}},
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
	svc := service{
		cardRegister:       mockCardGetter{},
		attachmentRegister: mockAttachmentGetter{},
		opener:             mockOpener{},
	}
	cardRegisterCalled = false
	attachmentRegisterCalled = false
	chargebackCreatorCalled = false

	_ = svc.create(disputeStub)

	if !cardRegisterCalled {
		t.Error("card register not called")
	}

	if !attachmentRegisterCalled {
		t.Error("attachment register not called")
	}

	if !chargebackCreatorCalled {
		t.Error("chargeback creator not called")
	}
}

func TestUnmarshalJSON_Errors(t *testing.T) {
	errDate := "unparseableData"
	cases := []struct {
		name string
		in   []byte
		want error
	}{
		{"null", []byte("null"), nil},
		{"parseError", []byte(errDate), &time.ParseError{
			Layout:     "2006-01-02",
			Value:      errDate,
			LayoutElem: "2006",
			ValueElem:  errDate,
			Message:    "",
		}}}
	var d date

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := d.UnmarshalJSON(c.in); !reflect.DeepEqual(got, c.want) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
