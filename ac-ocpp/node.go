package ac_ocpp

type ChargeStation struct {
	sn string
}

func (c *ChargeStation) ID() string {
	return c.sn
}
