package warehouse

import (
	"errors"
)

var ErrInvalidRating = errors.New("invalid rating - should be between 1 and 10 (including)")
var ErrCustomerNotAllowedToLeaveReview = errors.New("customer is not allowed to leave an error")

type DependencyOption func(w *Warehouse)

func WithPaymentProvider(paymentProvider PaymentProvider) DependencyOption {
	return func(w *Warehouse) {
		w.paymentProvider = paymentProvider
	}
}

func WithChartsProvider(chartsProvider ChartsProvider) DependencyOption {
	return func(w *Warehouse) {
		w.chartsProvider = chartsProvider
	}
}

type Warehouse struct {
	CDStock         []*CDBatch
	paymentProvider PaymentProvider
	chartsProvider  ChartsProvider
}

func NewWarehouse(opts ...DependencyOption) *Warehouse {
	w := &Warehouse{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func NewWarehouseWithPaymentProvider(provider PaymentProvider) *Warehouse {
	return &Warehouse{
		paymentProvider: provider,
	}
}

func (w *Warehouse) ReceiveBatchOfCDs(incomingBatches []CDBatch) {
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

func (w *Warehouse) AmountOfASpecificCDInStock(title, artist string) int {
	cdBatch := w.Search(title, artist)
	if cdBatch == nil {
		return 0
	}
	return cdBatch.Amount
}

func (w *Warehouse) Search(title, artist string) *CDBatch {
	lookingForCD := CD{
		Title:  title,
		Artist: artist,
	}

	for _, cdBatchInStock := range w.CDStock {
		if cdBatchInStock.CD.Equals(lookingForCD) {
			return cdBatchInStock
		}
	}

	return nil
}

func (w *Warehouse) SellCDToCustomer(cd *CD, customer *Customer) error {
	err := w.paymentProvider.ProcessPayment()
	if err != nil {
		return err
	}
	cdBatch := w.Search(cd.Title, cd.Artist)
	if cdBatch == nil {
		return errors.New("cd not found")
	}
	cdBatch.DecreaseAmount()
	customer.BuyCD(cd)
	return w.chartsProvider.Notify(cd.Title, cd.Artist, 1)
}

func (w *Warehouse) PriceForAlbum(title, artist string) int {
	position, competitorPrice := w.chartsProvider.PositionAndPriceForAlbum(title, artist)
	if position <= 100 {
		return competitorPrice - 1
	}
	return 100
}

func (w *Warehouse) LeaveReviewForCDByCustomer(cd *CD, review *Review, customer *Customer) error {
	if !customer.CanLeaveReviewForCD(cd) {
		return ErrCustomerNotAllowedToLeaveReview
	}
	if !review.IsValid() {
		return ErrInvalidRating
	}
	cdBatchInStock := w.Search(cd.Title, cd.Artist)
	if cdBatchInStock == nil {
		return nil
	}
	if cdBatchInStock.Reviews == nil {
		cdBatchInStock.Reviews = make([]*Review, 0)
	}
	cdBatchInStock.Reviews = append(cdBatchInStock.Reviews, review)
	return nil
}

func (w *Warehouse) GetReviewsForCD(title, artist string) []*Review {
	cdBatch := w.Search(title, artist)
	if cdBatch == nil {
		return nil
	}

	return cdBatch.Reviews
}

func (w *Warehouse) addBatchToStock(cdBatch CDBatch) {
	w.CDStock = append(w.CDStock, &cdBatch)
}
