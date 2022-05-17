package warehouse

type Customer struct {
	boughtCDs map[CD]int
}

func (c *Customer) CanLeaveReviewForCD(cd *CD) bool {
	return c.HasBoughtCD(cd)
}

func (c *Customer) HasBoughtCD(cd *CD) bool {
	_, hasBought := c.boughtCDs[*cd]
	return hasBought
}

func (c *Customer) BuyCD(cd *CD) {
	if c.boughtCDs == nil {
		c.boughtCDs = make(map[CD]int)
	}
	if !c.HasBoughtCD(cd) {
		c.boughtCDs[*cd] = 1
		return
	}
	c.boughtCDs[*cd]++
}
