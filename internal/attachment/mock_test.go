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
	errUnsentFiles struct{}
	errSave        struct{}
	errArchiver    struct{}
)

func (e errStorageList) list(cid string, bucket string, path string) ([]protocol.File, error) {
	return nil, listError
}

func (e errStorageList) Get(cid string, bucket string, key string) (*protocol.File, error) {
	return nil, nil
}

func (e errStorageGet) list(cid string, bucket string, path string) ([]protocol.File, error) {
	return files, nil
}

func (e errStorageGet) Get(cid string, bucket string, key string) (*protocol.File, error) {
	return nil, getError
}

func (e errUnsentFiles) getUnsentFiles(*protocol.Dispute, []protocol.File) ([]protocol.File, error) {
	return nil, unsentFilesError
}

func (e errUnsentFiles) save(*protocol.Chargeback) error {
	return saveError
}

func (e errArchiver) compact(cid string, files []protocol.File, strToRemove string) (*protocol.Attachment, error) {
	return nil, archiverError
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
	m.listCalled = cid == ci && orgId == bucket && m.expPath == path

	return files, nil
}

func (m *mockStorage) Get(cid string, bucket string, key string) (*protocol.File, error) {
	if cid == cid && orgId == bucket && m.getFile(key) != nil {
		m.actGetCalled++
	}

	return m.getFile(key), nil
}

func (m *mockArchiver) compact(ci string, fs []protocol.File, strToRemove string) (*protocol.Attachment, error) {
	m.called = ci == cid && m.strToRemove == strToRemove && filesEquals(fs, getFiles)

	return attStub, nil
}

func (m *mockRepository) getUnsentFiles(d *protocol.Dispute, fs []protocol.File) ([]protocol.File, error) {
	m.getUnsentFilesCalled = d == disputeStub && filesEquals(fs, files)

	return unsentFiles, nil
}

func (m *mockRepository) save(c *protocol.Chargeback) error {
	m.saveCalled = c == chargebackStub

	return nil
}

func (e *errSave) getUnsentFiles(d *protocol.Dispute, fs []protocol.File) ([]protocol.File, error) {
	return unsentFiles, nil
}

func (e *errSave) save(c *protocol.Chargeback) error {
	return saveError
}

func filesEquals(files1 []protocol.File, files2 []protocol.File) (ok bool) {
	for i, f := range files1 {
		if ok = f.Key == files2[i].Key; ok {
			return
		}
	}
	return
}
