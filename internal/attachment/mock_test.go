package attachment

import "aws-poc/internal/protocol"

type (
	mockStorage struct {
		listCalled   bool
		expPath      string
		actGetCalled int
		expGetCalled int
		expFiles     [3][2]protocol.File
	}
	mockArchiver struct {
		called      bool
		strToRemove string
	}
	mockRepository struct {
		getUnsentFilesCalled bool
		saveCalled           bool
	}
	errStorageList struct{}
	errStorageGet  struct{}
	errAttachments struct{}
	errSave        struct{}
	errArchiver    struct{}
)

func (e errStorageList) list(cid string, bucket string, path string) ([]protocol.File, error) {
	return nil, errListStub
}

func (e errStorageList) Get(cid string, bucket string, key string) (*protocol.File, error) {
	return nil, nil
}

func (e errStorageGet) list(cid string, bucket string, path string) ([]protocol.File, error) {
	return files, nil
}

func (e errStorageGet) Get(cid string, bucket string, key string) (*protocol.File, error) {
	return nil, errGetStub
}

func (e errAttachments) Get(*protocol.Dispute, []protocol.File) ([]protocol.Attachment, error) {
	return nil, errAttachmentStub
}

func (e errAttachments) save(*protocol.Chargeback) error {
	return errSaveStub
}

func (e errArchiver) Compress(_ string, _ []protocol.File, _ string) ([]byte, error) {
	return nil, errArchiverStub
}

func (m *mockStorage) getCalled(expGetCalled int) bool {
	return expGetCalled == m.actGetCalled
}

func (m *mockStorage) getFile(key string) *protocol.File {
	for _, f := range m.expFiles {
		if ok := f[0].Key == key; ok {
			return &f[1]
		}
	}
	return nil
}

func (m *mockStorage) list(ci string, bucket string, path string) ([]protocol.File, error) {
	m.listCalled = cid == ci && orgID == bucket && m.expPath == path

	return files, nil
}

func (m *mockStorage) Get(cid string, bucket string, key string) (*protocol.File, error) {
	if cid == cid && orgID == bucket && m.getFile(key) != nil {
		m.actGetCalled++
	}

	return m.getFile(key), nil
}

func (m *mockArchiver) Compress(ci string, fs []protocol.File, strToRemove string) ([]byte, error) {
	m.called = ci == cid && m.strToRemove == strToRemove && filesEquals(fs, getFiles)

	return compactFilesStub, nil
}

func (m *mockRepository) Get(d *protocol.Dispute, fs []protocol.File) ([]protocol.Attachment, error) {
	m.getUnsentFilesCalled = d == disputeStub && filesEquals(fs, files)

	return attachments, nil
}

func (m *mockRepository) save(c *protocol.Chargeback) error {
	m.saveCalled = c == chargebackStub

	return nil
}

func (e *errSave) GetUnsentFiles(d *protocol.Dispute, fs []protocol.File) ([]protocol.File, error) {
	return unsentFiles, nil
}

func (e *errSave) save(c *protocol.Chargeback) error {
	return errSaveStub
}

func filesEquals(files1 []protocol.File, files2 []protocol.File) (ok bool) {
	for i, f := range files1 {
		if ok = f.Key == files2[i].Key; ok {
			return
		}
	}
	return
}
