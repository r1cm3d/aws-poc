package attachment

//
//import (
//	"aws-poc/internal/protocol"
//	"reflect"
//	"testing"
//)
//
//// FIXME: Adapt test for the new unsentFiles function
//func TestGetSuccess(t *testing.T) {
//	storage, archiver, repository, exp := storageStub, &mockArchiver{
//		strToRemove: path,
//	},
//		&mockRepository{},
//		attStub
//	svc := svc{
//		storage:    storage,
//		Compressor: archiver,
//		Repository: repository,
//	}
//	act, _ := svc.Get(disputeStub)
//
//	if !storage.listCalled {
//		t.Error("storage list not listCalled")
//	}
//	if !repository.getUnsentFilesCalled {
//		t.Error("repository Get unsent files not listCalled")
//	}
//	if !storage.getCalled(3) {
//		t.Errorf("storage Get not listCalled %d times", 3)
//	}
//	if !archiver.called {
//		t.Error("Compressor not listCalled")
//	}
//
//	if !reflect.DeepEqual(act, exp) {
//		t.Errorf("GetSuccess, want: %v, got: %v", act, exp)
//	}
//}
//
//// FIXME: the same above
//func TestGetFail(t *testing.T) {
//	cases := []struct {
//		name string
//		in   *protocol.Dispute
//		svc
//		want error
//	}{
//		{"errListStub", disputeStub, svc{
//			storage: errStorageList{},
//		}, errListStub},
//		{"errAttachmentStub", disputeStub, svc{
//			storage:    &mockStorage{},
//			Repository: &errAttachments{},
//		}, errAttachmentStub},
//		{"errGetStub", disputeStub, svc{
//			storage:    &errStorageGet{},
//			Repository: &mockRepository{},
//		}, errGetStub},
//		{"compactError", disputeStub, svc{
//			storage:    storageStub,
//			Repository: &mockRepository{},
//			Compressor: &errArchiver{},
//		}, errArchiverStub},
//	}
//
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			if _, got := c.svc.Get(c.in); got != c.want {
//				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
//			}
//		})
//	}
//}
//
//func TestSave(t *testing.T) {
//	cases := []struct {
//		name string
//		in   *protocol.Chargeback
//		svc
//		want error
//	}{
//		{"success", chargebackStub, svc{
//			Repository: &mockRepository{},
//		}, nil},
//		{"error", chargebackStub, svc{
//			Repository: &errAttachments{},
//		}, errSaveStub},
//	}
//
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//			if got := c.svc.Save(c.in); got != c.want {
//				t.Errorf("%s, want: %v, got: %v", c.name, c.want, got)
//			}
//		})
//	}
//}
