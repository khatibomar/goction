package item

import "sync"

// Item is simplest abstraction of an Item in an auction
type Item struct {
	mu            sync.Mutex
	Name          string
	StartingPrice uint64
	CurrentPrice  uint64
	Currency      string
}

// New will create a new Item with current price is same as starting price
func New(name string, startingPrice uint64, currency string) *Item {
	return &Item{
		Name:          name,
		StartingPrice: startingPrice,
		CurrentPrice:  startingPrice,
		Currency:      currency,
	}
}

// UpdatePrice will update the price of the item
// even if the new price is lower than old price
// the caller should decide the logic of controling the price
func (i *Item) UpdatePrice(newPrice uint64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.CurrentPrice = newPrice
}

// Lock will lock the struct
// The use case of this is when the auction on the Item ends
// TODO(khatibomar)[Urgent]: is it safe to never Release lock?
func (i *Item) Lock() {
	i.mu.Lock()
}
