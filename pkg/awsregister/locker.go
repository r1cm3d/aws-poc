package awsregister

import (
	"aws-poc/internal/dispute"
)

type locker struct {
	register
}

func (l locker) lock(dispute dispute.Entity) error {
	return l.put(dispute)
}
