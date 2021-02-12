package attachment

import "testing"

//TODO: implement it
func TestGetSuccess(t *testing.T) {
	storage, archiver, repository := mockStorage{}, mockArchiver{}, mockRepository{}
	svc := svc{
		storage: storage,
		archiver:  archiver,
		repository: repository,
	}
	_, att = svc.Get(disputeStub)

	if !storage.listCalled {
		t.Error("storage list not called")
	}
	if !repository.getUnsentFilesCalled {
		t.Error("repository get unsent files not called")
	}
	if !storage.getCalled(3) {
		t.Errorf("storage get not called %d times", 3)
	}
	if !archiver.called {
		t.Error("archiver not called")
	}

	if att != exp {
		t.Errorf("GetSuccess, want: %v, got: %v", att, exp)
	}
}
