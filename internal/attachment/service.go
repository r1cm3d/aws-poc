package attachment

type (
	Entity struct {
		Name   string
		Base64 string
	}

	request struct {
		cid       string
		disputeId int
		accountId int
		orgId     string
		files     []string
	}

	Register interface {
		get(request) (Entity, error)
		save(request) error
	}
)
