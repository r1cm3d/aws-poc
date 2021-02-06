package attachment

type (
	Entity struct {
		Name   string
		Base64 string
	}

	Request struct {
		Cid       string
		DisputeId int
		AccountId int
		OrgId     string
	}

	Register interface {
		Get(Request) (Entity, error)
		Save(Request) error
	}
)
