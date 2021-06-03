package service

func (c *ChargeStation) withSN(sn string) {
	c.entry.WithField("sn", sn)
}

func (t *Transaction) withID(id string) {
	t.entry.WithField("transactionID", id)
}
