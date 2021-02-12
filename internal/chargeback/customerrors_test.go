package chargeback

import (
	"fmt"
	"testing"
)

func TestErrorsMessages(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"parseErrorMessage", newParseError(stubError), fmt.Sprintf("parser error: %v", stubError.Error())},
		{"idempotenceErrorMessage", newIdempotenceError(cid, disputeID), fmt.Sprintf("idempotence error: cid(%v), disputeId(%v)", cid, disputeID)},
		{"chargebackErrorMessage", newChargebackError(stubError, cid, disputeID), fmt.Sprintf("chargeback error: src(%v), cid(%v), disputeId(%v), ", stubError, cid, disputeID)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.err.Error(); got != c.want {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
