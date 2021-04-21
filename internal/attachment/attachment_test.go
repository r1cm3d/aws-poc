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
		Archiver:   archiver,
		repository: repository,
	}
	act, _ := svc.Get(disputeStub)

	if !storage.listCalled {
		t.Error("storage list not listCalled")
	}
	if !repository.getUnsentFilesCalled {
		t.Error("repository Get unsent files not listCalled")
	}
	if !storage.getCalled(3) {
		t.Errorf("storage Get not listCalled %d times", 3)
	}
	if !archiver.called {
		t.Error("Archiver not listCalled")
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
		{"errListStub", disputeStub, svc{
			storage: errStorageList{},
		}, errListStub},
		{"errUnsentFilesStub", disputeStub, svc{
			storage:    &mockStorage{},
			repository: &errUnsentFiles{},
		}, errUnsentFilesStub},
		{"errGetStub", disputeStub, svc{
			storage:    &errStorageGet{},
			repository: &mockRepository{},
		}, errGetStub},
		{"compactError", disputeStub, svc{
			storage:    storageStub,
			repository: &mockRepository{},
			Archiver:   &errArchiver{},
		}, errArchiverStub},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, got := c.svc.Get(c.in); got != c.want {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}

func TestSave(t *testing.T) {
	cases := []struct {
		name string
		in   *protocol.Chargeback
		svc
		want error
	}{
		{"success", chargebackStub, svc{
			repository: &mockRepository{},
		}, nil},
		{"error", chargebackStub, svc{
			repository: &errUnsentFiles{},
		}, errSaveStub},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.svc.Save(c.in); got != c.want {
				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
			}
		})
	}
}
