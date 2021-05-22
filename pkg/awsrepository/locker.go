package awsrepository

import (
	"aws-poc/internal/protocol"
	"fmt"
)

type (
	lockerRepository struct {
		repository
	}
	lockerRecord struct {
		protocol.Dispute
	}
)

func (l lockerRecord) ID() string {
	return fmt.Sprintf("%d::%v", l.DisputeID, l.Cid)
}

func (l lockerRepository) lock(dispute protocol.Dispute) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.put(rec)

	return err == nil
}

func (l lockerRepository) release(dispute protocol.Dispute) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.delete(rec)

	return err == nil
}
