package warehouse

type Review struct {
	Rating int
	Text   string
}

func (r *Review) IsValid() bool {
	return r.Rating > 1 && r.Rating < 10
}
