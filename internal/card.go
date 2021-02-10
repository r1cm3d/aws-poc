package internal

type (
	Card struct {
		Number string
	}

	CardGetter interface {
		Get(dispute *Dispute) (*Card, error)
	}
)
