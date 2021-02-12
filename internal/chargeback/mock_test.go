package chargeback

import "aws-poc/internal/protocol"

type (
	errMapper          struct{}
	errRepository      struct{}
	errDisputer        struct{}
	mockRepository     struct{}
	mockMapper         struct{}
	mockCreator        struct{}
	mockCardRepository struct {
		called bool
	}
	mockAttachmentRepository struct {
		getCalled  bool
		saveCalled bool
	}
	mockOpener struct {
		called bool
	}
	mockProducer struct {
		called   bool
		expected *protocol.Chargeback
	}
	mockScheduler struct {
		called bool
	}
	errCardGetter               struct{}
	errAttachmentGetRepository  struct{}
	errAttachmentSaveRepository struct{}
	errOpener                   struct{}
	errProducer                 struct{}
	errScheduler                struct{}
	mockOpenerWithNetworkError  struct {
		called bool
	}
)

func (e errMapper) fromJSON(string, string) (*protocol.Dispute, error) {
	return disputeStub, stubError
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
	return stubError
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

func (m *mockCardRepository) Get(dispute *protocol.Dispute) (*protocol.Card, error) {
	m.called = dispute == disputeStub

	return cardStub, nil
}

func (m *mockAttachmentRepository) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	m.getCalled = dispute == disputeStub

	return attachmentStub, nil
}

func (m *mockAttachmentRepository) Save(chargeback *protocol.Chargeback) error {
	m.saveCalled = chargeback == chargebackStub

	return nil
}

func (m *mockOpener) Open(dispute *protocol.Dispute, card *protocol.Card, attachment *protocol.Attachment) (*protocol.Chargeback, error) {
	m.called = dispute == disputeStub && card == cardStub && attachment == attachmentStub

	return chargebackStub, nil
}

func (m errCardGetter) Get(dispute *protocol.Dispute) (*protocol.Card, error) {
	return nil, cardError
}

func (m errAttachmentGetRepository) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	return nil, attGetError
}

func (m errAttachmentGetRepository) Save(chargeback *protocol.Chargeback) error {
	return nil
}

func (m errAttachmentSaveRepository) Get(dispute *protocol.Dispute) (*protocol.Attachment, error) {
	return nil, nil
}

func (m errAttachmentSaveRepository) Save(chargeback *protocol.Chargeback) error {
	return attSaveError
}

func (m errOpener) Open(dispute *protocol.Dispute, card *protocol.Card, attachment *protocol.Attachment) (*protocol.Chargeback, error) {
	return nil, openerError
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
	return producerError
}

func (e errScheduler) Schedule(chargeback *protocol.Chargeback) error {
	return scdError
}

func (m *mockOpenerWithNetworkError) Open(dispute *protocol.Dispute, card *protocol.Card, attachment *protocol.Attachment) (*protocol.Chargeback, error) {
	m.called = dispute == disputeStub && card == cardStub && attachment == attachmentStub

	return chargebackWithErrorStub, nil
}
