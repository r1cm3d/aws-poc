package awsregister

import (
	"aws-poc/internal/dispute"
	"fmt"
)

type locker struct {
	register
}

type lockerRecord struct {
	dispute.Entity
}

func (l lockerRecord) ID() string {
	return fmt.Sprintf("%s::%d", l.CorrelationID, l.DisputeID)
}

func (l locker) lock(dispute dispute.Entity) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.put(rec)

	return err == nil
}

func (l locker) release(dispute dispute.Entity) (ok bool) {
	rec := lockerRecord{dispute}
	err := l.delete(rec)

	return err == nil
}
