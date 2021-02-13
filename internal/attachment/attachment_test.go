package attachment

import (
	"testing"
)

func TestGetSuccess(t *testing.T) {
	storage, archiver, repository, exp := &mockStorage{
		expPath:  path,
		expFiles: [3][2]file{{uf1, fg1}, {uf2, fg2}, {uf3, fg3}},
	}, &mockArchiver{
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
