package item_test

import (
	"testing"

	"github.com/khatibomar/goction/item"
)

func TestUpdatePrice(t *testing.T) {
	testCases := []struct {
		name     string
		i        *item.Item
		newPrice uint64
		expPrice uint64
		expErr   error
	}{
		{"valid", item.New("usd item", 10, "$"), 20, 20, nil},
		{"same price", item.New("euro same price", 10, "euro"), 10, 10, nil},
		{"LBP valid", item.New("LBP item", 20000, "LBP"), 50000, 50000, nil},
		{"lower price", item.New("lower LBP item", 20000, "LBP"), 10000, 20000, item.ErrLower},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.i.UpdatePrice(tc.newPrice)
			if err != tc.expErr {
				t.Errorf("excpected %s , got %s", tc.expErr, err)
			}
			if tc.i.GetPrice() != tc.expPrice {
				t.Errorf("excpected %d , got %d", tc.newPrice, tc.expPrice)
			}
		})
	}
}

func TestLockedItem(t *testing.T) {
	i := item.New("test", 10, "$")
	i.Lock()
	err := i.UpdatePrice(20)
	if err != item.ErrLocked {
		t.Errorf("excpected %s , got %s", item.ErrLocked, err)
	}
}
