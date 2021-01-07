package dispute

import (
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

func (e errMapper) fromJSON(string, string) (Entity, error) {
	return disputeFake, errFake
}

func (m mockMapper) fromJSON(string, string) (Entity, error) {
	return disputeFake, nil
}

func (e errRepository) lock(Entity) (ok bool) {
	return false
}

func (e errRepository) release(Entity) (ok bool) {
	return false
}

func (e errDisputer) open(Entity) error {
	return errFake
}

func (m mockDisputer) open(Entity) error {
	return nil
}

func (m mockRepository) lock(Entity) (ok bool) {
	return true
}

func (m mockRepository) release(Entity) (ok bool) {
	return true
}

func TestMapFromJson(t *testing.T) {
	svc := service{}
	cid := cid
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
	want := Entity{
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
	defaultInput := [2]string{cid, "body"}
	cases := []struct {
		name string
		in   [2]string
		want error
		service
	}{
		{"success", defaultInput, nil, service{mapper: mockMapper{}, locker: mockRepository{}, disputer: mockDisputer{}}},
		{"parseError", defaultInput, newParseError(errFake), service{locker: mockRepository{}, mapper: errMapper{}}},
		{"idempotenceError", defaultInput, newIdempotenceError(cid, disputeID), service{mapper: mockMapper{}, locker: errRepository{}}},
		{"chargebackError", defaultInput, newChargebackError(errFake, cid, disputeID), service{mapper: mockMapper{}, locker: mockRepository{}, disputer: errDisputer{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.service.handleMessage(c.in[0], c.in[1]); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	svc := service{}
	d := Entity{}

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
