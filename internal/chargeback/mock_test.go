package chargeback

import "aws-poc/internal/protocol"

type (
	errMapper       struct{}
	errRepository   struct{}
	errDisputer     struct{}
	mockRepository  struct{}
	mockMapper      struct{}
	mockCreator     struct{}
	mockCardService struct {
		called bool
	}
	mockAttService struct {
		getCalled  bool
		saveCalled bool
	}
	mockNetworkCreator struct {
		called bool
	}
	mockProducer struct {
		called   bool
		expected *protocol.Chargeback
	}
	mockScheduler struct {
		called bool
	}
	errCardService             struct{}
	errAttGetService           struct{}
	errAttSaveService          struct{}
	errNetworkCreator          struct{}
	errProducer                struct{}
	errScheduler               struct{}
	mockOpenerWithNetworkError struct {
		called bool
	}
)

func (e errMapper) fromJSON(string, string) (*protocol.Dispute, error) {
	return disputeStub, errStub
}

func (m mockMapper) fromJSON(string, string) (*protocol.Dispute, error) {
	return disputeStub, nil
}

func (e errRepository) lock(*protocol.Dispute) (ok bool) {
	return false
}

func (e errRepository) release(*protocol.Dispute) (ok bool) {
	return false
}

func (e errDisputer) create(*protocol.Dispute) error {
	return errStub
}

func (m mockCreator) create(*protocol.Dispute) error {
	return nil
}

func (m mockRepository) lock(*protocol.Dispute) (ok bool) {
	return true
}

func (m mockRepository) release(*protocol.Dispute) (ok bool) {
	return true
}

func (m *mockCardService) Get(dispute *protocol.Dispute) (*protocol.Card, error) {
	m.called = dispute == disputeStub

	return cardStub, nil
}

func (m *mockAttService) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	m.getCalled = dispute == disputeStub

	return attachmentStub, nil
}

func (m *mockAttService) Save(chargeback *protocol.Chargeback) error {
	m.saveCalled = chargeback == chargebackStub

	return nil
}

func (m *mockNetworkCreator) Create(dispute *protocol.Dispute, card *protocol.Card, attachment *protocol.Attachment) (*protocol.Chargeback, error) {
	m.called = dispute == disputeStub && card == cardStub && attachment == attachmentStub

	return chargebackStub, nil
}

func (m errCardService) Get(dispute *protocol.Dispute) (*protocol.Card, error) {
	return nil, errCardStub
}

func (m errAttGetService) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	return nil, errAttGetStub
}

func (m errAttGetService) Save(chargeback *protocol.Chargeback) error {
	return nil
}

func (m errAttSaveService) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	return nil, nil
}

func (m errAttSaveService) Save(chargeback *protocol.Chargeback) error {
	return errAttSaveStub
}

func (m errNetworkCreator) Create(dispute *protocol.Dispute, card *protocol.Card, attachment *protocol.Attachment) (*protocol.Chargeback, error) {
	return nil, errOpenerStub
}

func (m *mockProducer) Produce(chargeback *protocol.Chargeback) error {
	m.called = m.expected == chargeback

	return nil
}

func (m *mockScheduler) Schedule(chargeback *protocol.Chargeback) error {
	m.called = chargeback == chargebackStub

	return nil
}

func (e errProducer) Produce(chargeback *protocol.Chargeback) error {
	return errProducerStub
}

func (e errScheduler) Schedule(chargeback *protocol.Chargeback) error {
	return errScdStub
}

func (m *mockOpenerWithNetworkError) Create(dispute *protocol.Dispute, card *protocol.Card, attachment *protocol.Attachment) (*protocol.Chargeback, error) {
	m.called = dispute == disputeStub && card == cardStub && attachment == attachmentStub

	return chargebackWithErrorStub, nil
}
