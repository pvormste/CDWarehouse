package warehouse

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

/*
- Warehouse gets batches of CDs (/)
- Warehouse has a stock of CDs (having count) (/)
- Search for CDs (/)
  - title
  - artist
- Customer can buy CDs
  - Payment by external provider
- CDs can have (/)
  - reviews from customers (rated 1 - 10)
    - optional: text
*/

type FakePaymentProvider struct {
	failPayment bool
}

func (f *FakePaymentProvider) ProcessPayment() error {
	if f.failPayment {
		return errors.New("payment failed")
	}
	return nil
}

func TestWarehouse(t *testing.T) {
	t.Run("warehouse gets batches of CDs", func(t *testing.T) {
		t.Run("empty batch", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.ReceiveBatchOfCDs([]CDBatch{})
			assert.Equal(t, 0, warehouse.CDsInStock())
		})

		t.Run("sending a single CD", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.ReceiveBatchOfCDs([]CDBatch{
				{
					CD:     CD{},
					Amount: 1,
				},
			})
			assert.Equal(t, 1, warehouse.CDsInStock())
		})

		t.Run("sending multiple CDs", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.ReceiveBatchOfCDs([]CDBatch{
				{
					CD: CD{
						Title:  "Viva la Vida",
						Artist: "Coldplay",
					},
					Amount: 2,
				},
				{
					CD: CD{
						Title:  "Amerika",
						Artist: "Rammstein",
					},
					Amount: 1,
				},
			})
			assert.Equal(t, 3, warehouse.CDsInStock())
		})
	})

	t.Run("search for CDs in warehouse", func(t *testing.T) {
		t.Run("CD can't be found by title and artist", func(t *testing.T) {
			warehouse := NewWarehouse()
			foundCDBatch := warehouse.Search("Amerika", "Rammstein")
			assert.Nil(t, foundCDBatch)
		})

		t.Run("CD can't be found because it doesn't exist with this title", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.ReceiveBatchOfCDs([]CDBatch{
				{
					CD: CD{
						Title:  "Amerika",
						Artist: "Rammstein",
					},
					Amount: 1,
				},
			})
			foundCDBatch := warehouse.Search("Rosenrot", "Rammstein")
			assert.Nil(t, foundCDBatch)
		})

		t.Run("CD can be found by title and artist", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.ReceiveBatchOfCDs([]CDBatch{
				{
					CD: CD{
						Title:  "Amerika",
						Artist: "Rammstein",
					},
					Amount: 2,
				},
			})
			foundCDBatch := warehouse.Search("Amerika", "Rammstein")
			expectedCDBatch := &CDBatch{
				CD: CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				},
				Amount: 2,
			}
			assert.Equal(t, expectedCDBatch, foundCDBatch)
		})
	})

	t.Run("CD reviews", func(t *testing.T) {
		t.Run("customer cannot leave review without buying", func(t *testing.T) {
			customer := &Customer{
				boughtCDs: nil,
			}
			cd := &CD{
				Title:  "Amerika",
				Artist: "Rammstein",
			}
			assert.False(t, customer.CanLeaveReviewForCD(cd))
		})

		t.Run("customer can leave review if bought CD", func(t *testing.T) {
			cd := CD{
				Title:  "Amerika",
				Artist: "Rammstein",
			}
			customer := &Customer{
				boughtCDs: map[CD]int{
					cd: 1,
				},
			}
			assert.True(t, customer.CanLeaveReviewForCD(&cd))
		})

		t.Run("customer can leave review with rating only", func(t *testing.T) {
			t.Run("will return error if rating is below 1", func(t *testing.T) {
				warehouse := NewWarehouse()
				cd := CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				}
				review := Review{
					Rating: 0,
					Text:   "",
				}
				customer := Customer{
					boughtCDs: map[CD]int{
						cd: 1,
					},
				}
				err := warehouse.LeaveReviewForCDByCustomer(&cd, &review, &customer)
				assert.Error(t, err)
			})

			t.Run("will return error if rating is higher than 10", func(t *testing.T) {
				warehouse := NewWarehouse()
				cd := CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				}
				review := Review{
					Rating: 11,
					Text:   "",
				}
				customer := Customer{
					boughtCDs: map[CD]int{
						cd: 1,
					},
				}
				err := warehouse.LeaveReviewForCDByCustomer(&cd, &review, &customer)
				assert.Error(t, err)
			})

			t.Run("will be successfully adding a review to the CD", func(t *testing.T) {
				warehouse := NewWarehouse()
				cd := CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				}
				warehouse.ReceiveBatchOfCDs([]CDBatch{
					{
						CD:     cd,
						Amount: 1,
					},
				})
				review := Review{
					Rating: 5,
					Text:   "",
				}
				customer := Customer{boughtCDs: map[CD]int{
					cd: 1,
				}}
				err := warehouse.LeaveReviewForCDByCustomer(&cd, &review, &customer)
				actualReviews := warehouse.GetReviewsForCD("Amerika", "Rammstein")
				expectedReviews := []*Review{
					&review,
				}
				assert.NoError(t, err)
				assert.Equal(t, expectedReviews, actualReviews)
			})
			t.Run("customer who hasn't bought a CD can't leave a review", func(t *testing.T) {
				warehouse := NewWarehouse()
				cd := CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				}
				warehouse.ReceiveBatchOfCDs([]CDBatch{
					{
						CD:     cd,
						Amount: 1,
					},
				})
				review := Review{
					Rating: 8,
					Text:   "",
				}
				customer := Customer{}
				err := warehouse.LeaveReviewForCDByCustomer(&cd, &review, &customer)
				actualReviews := warehouse.GetReviewsForCD("Amerika", "Rammstein")
				assert.Error(t, err)
				assert.Equal(t, 0, len(actualReviews))
			})
		})

		t.Run("warehouse selling a CD to a customer", func(t *testing.T) {
			t.Run("payment doesn't work", func(t *testing.T) {
				warehouse := NewWarehouseWithPaymentProvider(&FakePaymentProvider{failPayment: true})
				cd := CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				}
				warehouse.ReceiveBatchOfCDs([]CDBatch{
					{
						CD:     cd,
						Amount: 5,
					},
				})

				customer := Customer{}

				err := warehouse.SellCDToCustomer(&cd, &customer)
				assert.Error(t, err)
				assert.Equal(t, 5, warehouse.AmountOfASpecificCDInStock("Amerika", "Rammstein"))
				assert.False(t, customer.HasBoughtCD(&cd))
			})

			t.Run("payment does work and reduces the stock", func(t *testing.T) {
				warehouse := NewWarehouseWithPaymentProvider(&FakePaymentProvider{})
				cd := CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				}
				warehouse.ReceiveBatchOfCDs([]CDBatch{
					{
						CD:     cd,
						Amount: 5,
					},
				})

				customer := Customer{}

				err := warehouse.SellCDToCustomer(&cd, &customer)
				assert.NoError(t, err)
				assert.Equal(t, 4, warehouse.AmountOfASpecificCDInStock("Amerika", "Rammstein"))
				assert.True(t, customer.HasBoughtCD(&cd))
			})

			t.Run("should notify charts on payment", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				chartsProviderMock := NewMockChartsProvider(ctrl)
				chartsProviderMock.EXPECT().
					Notify("Amerika", "Rammstein", 1).
					Return(nil)

				warehouse := NewWarehouse(
					WithPaymentProvider(&FakePaymentProvider{}),
					WithChartsProvider(chartsProviderMock),
				)
				cd := CD{
					Title:  "Amerika",
					Artist: "Rammstein",
				}
				warehouse.ReceiveBatchOfCDs([]CDBatch{
					{
						CD:     cd,
						Amount: 5,
					},
				})
				customer := Customer{}
				err := warehouse.SellCDToCustomer(&cd, &customer)
				assert.NoError(t, err)
			})
		})

		t.Run("Price for Album", func(t *testing.T) {
			t.Run("Get normal price for an Album if its not in the Top 100", func(t *testing.T) {
				mockCtrl := gomock.NewController(t)
				mockChartsProvider := NewMockChartsProvider(mockCtrl)
				mockChartsProvider.EXPECT().
					PositionAndPriceForAlbum("Amerika", "Rammstein").
					Return(150, 95)

				warehouse := NewWarehouse(
					WithChartsProvider(mockChartsProvider),
				)
				warehouse.ReceiveBatchOfCDs([]CDBatch{
					{
						CD: CD{
							Title:  "Amerika",
							Artist: "Rammstein",
						},
						Amount:  5,
						Reviews: nil,
					},
				})

				albumPrice := warehouse.PriceForAlbum("Amerika", "Rammstein")
				assert.Equal(t, 100, albumPrice)
			})

			t.Run("Get a 1 pound lower price if the album is in the top 100 charts", func(t *testing.T) {
				mockCtrl := gomock.NewController(t)
				mockChartsProvider := NewMockChartsProvider(mockCtrl)
				mockChartsProvider.EXPECT().
					PositionAndPriceForAlbum("Wait For U", "Future Featuring Drake & Tems").
					Return(1, 80)

				warehouse := NewWarehouse(
					WithChartsProvider(mockChartsProvider),
				)
				warehouse.ReceiveBatchOfCDs([]CDBatch{
					{
						CD: CD{
							Title:  "Wait For U",
							Artist: "Future Featuring Drake & Tems",
						},
						Amount:  5,
						Reviews: nil,
					},
				})

				albumPrice := warehouse.PriceForAlbum("Wait For U", "Future Featuring Drake & Tems")
				assert.Equal(t, 79, albumPrice)
			})
		})
	})
}
