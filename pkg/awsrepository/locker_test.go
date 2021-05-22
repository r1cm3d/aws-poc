package awsrepository

import (
	"aws-poc/internal/protocol"
	"aws-poc/pkg/awssession"
	"aws-poc/pkg/test/integration"
	"reflect"
	"testing"
)

func TestLockIntegration(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()
	cases := []struct {
		name string
		in   protocol.Dispute
		want bool
		lockerRepository
	}{
		{"success", *disputeStub, true, lockerRepository{newRegister(awssession.NewLocalSession(), tableName)}},
		{"error", *disputeStub, false, lockerRepository{errRegister{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.lockerRepository.lock(c.in); !reflect.DeepEqual(c.want, got) {
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
		lockerRepository
	}{
		{"success", *disputeStub, true, lockerRepository{newRegister(awssession.NewLocalSession(), tableName)}},
		{"error", *disputeStub, false, lockerRepository{errRegister{}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.lockerRepository.release(c.in); !reflect.DeepEqual(c.want, got) {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
