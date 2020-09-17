package dispute

import "errors"

const (
	disputeID = 666
	cid       = "e1388e36-1683-4902-b30c-5c5b63f5976c"
)

var (
	errFake     = errors.New("mocked error")
	disputeFake = dispute{DisputeID: disputeID}
)
