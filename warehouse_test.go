package warehouse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
- Warehouse gets batches of CDs
- Warehouse has a stock of CDs (having count)
- Search for CDs
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
}
