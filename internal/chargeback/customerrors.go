package chargeback

import "fmt"

type (
	parseError struct {
		src error
	}
	idempotenceError struct {
		cid       string
		disputeID int
	}
	chargebackError struct {
		src       error
		cid       string
		disputeID int
	}
)

func newParseError(src error) error {
	return &parseError{src}
}

func newIdempotenceError(cid string, disputeID int) error {
	return &idempotenceError{
		cid,
		disputeID,
	}
}

func newChargebackError(src error, cid string, disputeID int) error {
	return &chargebackError{
		src,
		cid,
		disputeID,
	}
}

func (e parseError) Error() string {
	return fmt.Sprintf("parser error: %s", e.src.Error())
}

func (e idempotenceError) Error() string {
	return fmt.Sprintf("idempotence error: cid(%v), disputeId(%v)", e.cid, e.disputeID)
}

func (e chargebackError) Error() string {
	return fmt.Sprintf("chargeback error: src(%v), cid(%v), disputeId(%v), ", e.src, e.cid, e.disputeID)
}
