package warehouse

type CDBatch struct {
	CD      CD
	Amount  int
	Reviews []*Review
}

func (c *CDBatch) DecreaseAmount() {
	if c.Amount > 0 {
		c.Amount--
	}
}
