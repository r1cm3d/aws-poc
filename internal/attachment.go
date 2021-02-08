package internal

type (
	Attachment struct {
		Name   string
		Base64 string
	}

	AttachmentGetter interface {
		Get(dispute Dispute) (Attachment, error)
	}
)
