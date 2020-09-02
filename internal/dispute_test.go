package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestMapFromJson(t *testing.T) {
	svc := disputeSvc{}
	json := `{
  "correlationId": "7658c09d-a8c3-47f4-b584-922641ab3416",
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
	want := Dispute{
		CorrelationId:       "7658c09d-a8c3-47f4-b584-922641ab3416",
		DisputeId:           611,
		AccountId:           48448,
		AuthorizationCode:   "451",
		ReasonCode:          "848",
		CardId:              "3123",
		Tenant:              "pismo.io",
		DisputeAmount:       32.32,
		TransactionDate:     Date(time.Date(2012, 04, 23, 0, 0, 0, 0, time.UTC)),
		LocalCurrencyCode:   "USD",
		TextMessage:         "this is a test message",
		DocumentIndicator:   true,
		IsPartialChargeback: false,
	}

	got, err := svc.mapFromJson(json)

	if err != nil {
		t.Error("mapFromJson() error should not be returned")
	}
	if reflect.DeepEqual(got, want) {
		t.Errorf("mapFromJson() got: %v, want: %v", got, want)
	}
}

func TestMapFromJson_Error(t *testing.T) {
	svc := disputeSvc{}

	_, err := svc.mapFromJson("json")

	if err == nil {
		t.Error("mapFromJson_Error() error should be returned")
	}
}

func TestOpenChargeback(t *testing.T) {
	svc := disputeSvc{}
	d := Dispute{}

	if err := svc.openChargeback(d); err != nil {
		t.Error("openChargeback error should be returned")
	}
}
