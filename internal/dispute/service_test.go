package dispute

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type (
	errMapper      struct{}
	errRepository  struct{}
	errDisputer    struct{}
	mockRepository struct{}
	mockMapper     struct{}
	mockDisputer   struct{}
)

var errFake = errors.New("mocked error")

func (e errMapper) fromJSON(string, string) (dispute, error) {
	return dispute{}, errFake
}

func (m mockMapper) fromJSON(string, string) (dispute, error) {
	return dispute{}, nil
}

func (e errRepository) lock(dispute) (ok bool) {
	return false
}

func (e errRepository) unlock(dispute) {
}

func (e errDisputer) open(dispute) error {
	return errFake
}

func (m mockDisputer) open(dispute) error {
	return nil
}

func (m mockRepository) lock(dispute) (ok bool) {
	return true
}

func (m mockRepository) unlock(dispute) {
}

func TestMapFromJson(t *testing.T) {
	svc := service{}
	cid := "7658c09d-a8c3-47f4-b584-922641ab3416"
	json := `{
  "disputeId": 611,
  "accountId": 48448,
  "authorizationCode": "451",
  "reasonCode": "848",
  "cardId": "3123",
  "tenant": "pismo.io",
  "disputeAmount": 32.32,
  "transactionAmount": 42.65,
  "transactionDate": "2012-04-23",
  "localCurrencyCode": "USD",
  "textMessage": "this a test message",
  "documentIndicator": true,
  "isPartialChargeback": false
}`
	want := dispute{
		CorrelationID:       cid,
		DisputeID:           611,
		AccountID:           48448,
		AuthorizationCode:   "451",
		ReasonCode:          "848",
		CardID:              "3123",
		Tenant:              "pismo.io",
		DisputeAmount:       32.32,
		TransactionDate:     date(time.Date(2012, 04, 23, 0, 0, 0, 0, time.UTC)),
		LocalCurrencyCode:   "USD",
		TextMessage:         "this is a test message",
		DocumentIndicator:   true,
		IsPartialChargeback: false,
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
	svc := service{
		register: mockRepository{},
		mapper:   mockMapper{},
		disputer: mockDisputer{},
	}

	if err := svc.handleMessage("", ""); err != nil {
		t.Errorf("handleMessage should not return an error")
	}
}

func TestHandleMessage_Error(t *testing.T) {
	cases := []struct {
		name string
		service
	}{
		{"lockError", service{mapper: mockMapper{}, register: errRepository{}}},
		{"mapError", service{register: mockRepository{}, mapper: errMapper{}}},
		{"openError", service{mapper: mockMapper{}, register: mockRepository{}, disputer: errDisputer{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.service.handleMessage("", ""); err == nil {
				t.Errorf("%s . An error should be returned", c.name)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	svc := service{}
	d := dispute{}

	if err := svc.open(d); err != nil {
		t.Error("open error should be returned")
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