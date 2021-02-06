package card

type (
	Entity struct {
		Number string
	}

	Request struct {
		Cid       string
		OrgId     string
		AccountId int
	}

	Register interface {
		Get(Request) (Entity, error)
	}
)
