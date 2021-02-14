package attachment

import "aws-poc/internal/protocol"

type (
	mockStorage struct {
		listCalled   bool
		expPath      string
		actGetCalled int
		expGetCalled int
		expFiles     [3][2]file
	}
	mockArchiver struct {
		called      bool
		strToRemove string
	}
	mockRepository struct {
		getUnsentFilesCalled bool
	}
	errStorageList struct{}
	errStorageGet  struct{}
	errRepository  struct{}
	errArchiver    struct{}
)

func (e errStorageList) list(cid string, bucket string, path string) ([]file, error) {
	return nil, listError
}

func (e errStorageList) get(cid string, bucket string, key string) (*file, error) {
	return nil, nil
}

func (e errStorageGet) list(cid string, bucket string, path string) ([]file, error) {
	return files, nil
}

func (e errStorageGet) get(cid string, bucket string, key string) (*file, error) {
	return nil, getError
}

func (e errRepository) getUnsentFiles(*protocol.Dispute, []file) ([]file, error) {
	return nil, unsentFilesError
}

func (e errArchiver) compact(cid string, files []file, strToRemove string) (*protocol.Attachment, error) {
	return nil, archiverError
}

func (m *mockStorage) getCalled(expGetCalled int) bool {
	return expGetCalled == m.actGetCalled
}

func (m *mockStorage) getFile(key string) *file {
	for _, f := range m.expFiles {
		if ok := f[0].key == key; ok {
			return &f[1]
		}
	}
	return nil
}

func (m *mockStorage) list(ci string, bucket string, path string) ([]file, error) {
	m.listCalled = cid == ci && orgId == bucket && m.expPath == path

	return files, nil
}

func (m *mockStorage) get(cid string, bucket string, key string) (*file, error) {
	if cid == cid && orgId == bucket && m.getFile(key) != nil {
		m.actGetCalled++
	}

	return m.getFile(key), nil
}

func (m *mockArchiver) compact(ci string, fs []file, strToRemove string) (*protocol.Attachment, error) {
	m.called = ci == cid && m.strToRemove == strToRemove && filesEquals(fs, getFiles)

	return attStub, nil
}

func (m *mockRepository) getUnsentFiles(d *protocol.Dispute, fs []file) ([]file, error) {
	m.getUnsentFilesCalled = d == disputeStub && filesEquals(fs, files)

	return unsentFiles, nil
}

func filesEquals(files1 []file, files2 []file) (ok bool) {
	for i, f := range files1 {
		if ok = f.key == files2[i].key; ok {
			return
		}
	}
	return
}
