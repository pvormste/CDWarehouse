package warehouse

import (
	"testing"

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
- CDs can have
  - reviews from customers (rated 1 - 10)
    - optional: text
*/

func TestWarehouse(t *testing.T) {
	t.Run("warehouse gets batches of CDs", func(t *testing.T) {
		t.Run("empty batch", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.SendBatchOfCDs([]CDBatch{})
			assert.Equal(t, 0, warehouse.CDsInStock())
		})

		t.Run("sending a single CD", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.SendBatchOfCDs([]CDBatch{
				{
					CD:     CD{},
					Amount: 1,
				},
			})
			assert.Equal(t, 1, warehouse.CDsInStock())
		})

		t.Run("sending multiple CDs", func(t *testing.T) {
			warehouse := NewWarehouse()
			warehouse.SendBatchOfCDs([]CDBatch{
				{
					CD: CD{
						Artist: "Coldplay",
					},
					Amount: 2,
				},
				{
					CD: CD{
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
			warehouse.SendBatchOfCDs([]CDBatch{
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
			warehouse.SendBatchOfCDs([]CDBatch{
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
			warehouse := Warehouse{}
			customer := &Customer{
				BoughtCDs: nil,
			}
			cd := &CD{
				Title:  "Amerika",
				Artist: "Rammstein",
			}
			assert.False(t, warehouse.CustomerCanLeaveReviewForCD(customer, cd))
		})

		t.Run("customer can leave review if bought CD", func(t *testing.T) {
			warehouse := Warehouse{}
			cd := CD{
				Title:  "Amerika",
				Artist: "Rammstein",
			}
			customer := &Customer{
				BoughtCDs: map[CD]int{
					cd: 1,
				},
			}
			assert.True(t, warehouse.CustomerCanLeaveReviewForCD(customer, &cd))
		})
	})
}
