package card

type (
	Entity struct {
		Number string
	}

	request struct {
		cid       string
		orgId     string
		accountId int
	}

	Register interface {
		get(request) (Entity, error)
	}
)
