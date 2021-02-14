package attachment

import (
	"aws-poc/internal/protocol"
	"testing"
)

func TestGetSuccess(t *testing.T) {
	storage, archiver, repository, exp := storageStub, &mockArchiver{
		strToRemove: path,
	},
		&mockRepository{},
		attStub
	svc := svc{
		storage:    storage,
		archiver:   archiver,
		repository: repository,
	}
	act, _ := svc.Get(disputeStub)

	if !storage.listCalled {
		t.Error("storage list not listCalled")
	}
	if !repository.getUnsentFilesCalled {
		t.Error("repository get unsent files not listCalled")
	}
	if !storage.getCalled(3) {
		t.Errorf("storage get not listCalled %d times", 3)
	}
	if !archiver.called {
		t.Error("archiver not listCalled")
	}

	if act != exp {
		t.Errorf("GetSuccess, want: %v, got: %v", act, exp)
	}
}

func TestGetFail(t *testing.T) {
	cases := []struct {
		name string
		in   *protocol.Dispute
		svc
		want error
	}{
		{"listError", disputeStub, svc{
			storage: errStorageList{},
		}, listError},
		{"unsentFilesError", disputeStub, svc{
			storage:    &mockStorage{},
			repository: &errRepository{},
		}, unsentFilesError},
		{"getError", disputeStub, svc{
			storage:    &errStorageGet{},
			repository: &mockRepository{},
		}, getError},
		{"compactError", disputeStub, svc{
			storage:    storageStub,
			repository: &mockRepository{},
			archiver:   &errArchiver{},
		}, archiverError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, got := c.svc.Get(c.in); got != c.want {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
