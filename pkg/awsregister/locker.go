package awsregister

import (
	"aws-poc/internal"
	"fmt"
)

type (
	locker struct {
		register
	}
	lockerRecord struct {
		internal.Dispute
	}
)

func (l lockerRecord) ID() string {
	return fmt.Sprintf("%d::%v", l.DisputeId, l.Cid)
}

func (l locker) lock(dispute internal.Dispute) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.put(rec)

	return err == nil
}

func (l locker) release(dispute internal.Dispute) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.delete(rec)

	return err == nil
}
