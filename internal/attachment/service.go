package attachment

type (
	Entity struct {
		Name   string
		Base64 string
	}

	Register interface {
		Get(cid string, orgId string, accountId int, disputeId int) (Entity, error)
	}
)
