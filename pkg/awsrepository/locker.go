package awsrepository

import (
	"aws-poc/internal/protocol"
	"fmt"
)

type (
	locker struct {
		repository
	}
	lockerRecord struct {
		protocol.Dispute
	}
)

func (l lockerRecord) ID() string {
	return fmt.Sprintf("%d::%v", l.DisputeID, l.Cid)
}

func (l locker) lock(dispute protocol.Dispute) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.put(rec)

	return err == nil
}

func (l locker) release(dispute protocol.Dispute) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.delete(rec)

	return err == nil
}
