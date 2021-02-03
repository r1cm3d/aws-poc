package chargeback

const (
	CHARGEBACK         = Type("CHARGEBACK")
	SECOND_PRESENTMENT = Type("SECOND_PRESENTMENT")
	ARB_CHARGEBACK     = Type("ARB_CHARGEBACK")
)

type (
	Entity struct {
	}

	Type string

	request struct {
		cid               string
		currencyCode      string
		documentIndicator bool
		message           string
		disputedAmount    float64
		reasonCode        string
		isPartial         bool
		chargebackType    Type
	}

	Creator interface {
		create(request) (Entity, error)
	}

	Producer interface {
		produce(chargeback Entity) error
	}

	Register interface {
		save(chargeback Entity) error
	}

	Scheduler interface {
		scheduleForTomorrow(chargeback Entity) error
	}

	Facade interface {
		Creator
		Producer
		Register
		Scheduler
	}
)
