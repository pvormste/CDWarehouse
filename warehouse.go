package warehouse

type CD struct {
	Title  string
	Artist string
}

func (c *CD) Equals(otherCD CD) bool {
	return c.Title == otherCD.Title && c.Artist == otherCD.Artist
}

type CDBatch struct {
	CD     CD
	Amount int
}

type Warehouse struct {
	CDStock []CDBatch
}

func NewWarehouse() *Warehouse {
	return &Warehouse{}
}

func (w *Warehouse) SendBatchOfCDs(incomingBatches []CDBatch) {
	for _, incomingBatch := range incomingBatches {
		foundInStock := false
		for _, batchInStock := range w.CDStock {
			if batchInStock.CD.Equals(incomingBatch.CD) {
				batchInStock.Amount += incomingBatch.Amount
				foundInStock = true
			}
		}
		if !foundInStock {
			w.addBatchToStock(incomingBatch)
		}
	}

}

func (w *Warehouse) CDsInStock() int {
	cdCount := 0
	for _, batchInStock := range w.CDStock {
		cdCount += batchInStock.Amount
	}
	return cdCount
}

func (w *Warehouse) addBatchToStock(cdBatch CDBatch) {
	w.CDStock = append(w.CDStock, cdBatch)
}
