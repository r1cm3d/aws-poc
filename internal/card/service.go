package card

type (
	Entity struct {
		Number string
	}

	Register interface {
		Get(cid string, orgId string, accountId int) (Entity, error)
	}
)
