package warehouse

type CD struct {
	Title  string
	Artist string
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
			if batchInStock.CD.Title == incomingBatch.CD.Title && batchInStock.CD.Artist == incomingBatch.CD.Artist {
				batchInStock.Amount += incomingBatch.Amount
				foundInStock = true
			}
		}
		if !foundInStock {
			w.CDStock = append(w.CDStock, incomingBatch)
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
