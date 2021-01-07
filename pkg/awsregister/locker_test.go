package awsregister

import (
	"aws-poc/internal/dispute"
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"reflect"
	"testing"
)

type errRegister struct{}

var stubDispute = dispute.Entity{
	CorrelationID: "ee67f4f2-0b08-4f58-908f-bbb9bc37a1d2",
	DisputeID:     666,
}

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
		in   dispute.Entity
		want bool
		locker
	}{
		{"success", stubDispute, true, locker{newRegister(awssession.NewLocalSession(), tableName)}},
		{"error", stubDispute, false, locker{errRegister{}}},
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
		in   dispute.Entity
		want bool
		locker
	}{
		{"success", stubDispute, true, locker{newRegister(awssession.NewLocalSession(), tableName)}},
		{"error", stubDispute, false, locker{errRegister{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.locker.release(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
