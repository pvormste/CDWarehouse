package warehouse

type CD struct {
	Title  string
	Artist string
}

func (c *CD) Equals(otherCD CD) bool {
	return c.Title == otherCD.Title && c.Artist == otherCD.Artist
}
