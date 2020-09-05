package internal

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

func (e errMapper) mapFromJSON(string, string) (dispute, error) {
	return dispute{}, errFake
}

func (m mockMapper) mapFromJSON(string, string) (dispute, error) {
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
	svc := disputeSvc{}
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
		TransactionDate:     Date(time.Date(2012, 04, 23, 0, 0, 0, 0, time.UTC)),
		LocalCurrencyCode:   "USD",
		TextMessage:         "this is a test message",
		DocumentIndicator:   true,
		IsPartialChargeback: false,
	}

	got, err := svc.mapFromJSON(cid, json)

	if err != nil {
		t.Error("mapFromJSON() error should not be returned")
	}
	if reflect.DeepEqual(got, want) {
		t.Errorf("mapFromJSON() got: %v, want: %v", got, want)
	}
}

func TestMapFromJson_Error(t *testing.T) {
	svc := disputeSvc{}

	_, err := svc.mapFromJSON("", "json")

	if err == nil {
		t.Error("mapFromJson_Error() error should be returned")
	}
}

func TestHandleMessage(t *testing.T) {
	svc := disputeSvc{
		disputeRepository: mockRepository{},
		disputeMapper:     mockMapper{},
		disputer:          mockDisputer{},
	}

	if err := svc.handleMessage("", ""); err != nil {
		t.Errorf("handleMessage should not return an error")
	}
}

func TestHandleMessage_Error(t *testing.T) {
	cases := []struct {
		name string
		disputeSvc
	}{
		{"lockError", disputeSvc{disputeMapper: mockMapper{}, disputeRepository: errRepository{}}},
		{"mapError", disputeSvc{disputeRepository: mockRepository{}, disputeMapper: errMapper{}}},
		{"openError", disputeSvc{disputeMapper: mockMapper{}, disputeRepository: mockRepository{}, disputer: errDisputer{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.disputeSvc.handleMessage("", ""); err == nil {
				t.Errorf("%s . An error should be returned", c.name)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	svc := disputeSvc{}
	d := dispute{}

	if err := svc.open(d); err != nil {
		t.Error("open error should be returned")
	}
}
