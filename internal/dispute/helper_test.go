package dispute

import "errors"

const (
	disputeID = 666
	cid       = "e1388e36-1683-4902-b30c-5c5b63f5976c"
	orgId     = "TN-ed3d9cbf-664e-4044-bc1f-5adee7ff069f"
	accountId = 10782
)

var (
	errFake     = errors.New("mocked error")
	disputeFake = Entity{DisputeID: disputeID}
)
