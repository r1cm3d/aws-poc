package awsrepository

import (
	"aws-poc/internal/protocol"
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"reflect"
	"testing"
)

var disputeStub = protocol.Dispute{DisputeId: 123, Cid: "cid"}

type errRegister struct{}

func (e errRegister) put(_ record) error {
	return errPutItem
}

func (e errRegister) delete(_ record) error {
	return errDelete
}

func TestLockIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name string
		in   protocol.Dispute
		want bool
		locker
	}{
		{"success", disputeStub, true, locker{newRegister(awssession.NewLocalSession(), tableName)}},
		{"error", disputeStub, false, locker{errRegister{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.locker.lock(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestReleaseIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name string
		in   protocol.Dispute
		want bool
		locker
	}{
		{"success", disputeStub, true, locker{newRegister(awssession.NewLocalSession(), tableName)}},
		{"error", disputeStub, false, locker{errRegister{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.locker.release(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
