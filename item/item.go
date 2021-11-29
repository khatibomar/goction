package item

import (
	"errors"
	"sync"
)

var (
	ErrLocked = errors.New("Item is locked")
	ErrLower  = errors.New("New price is lower than old price")
)

// Item is simplest abstraction of an Item in an auction
type Item struct {
	mu            sync.RWMutex
	name          string
	startingPrice uint64
	currentPrice  uint64
	currency      string
	locked        bool
}

// New will create a new Item with current price is same as starting price
func New(name string, startingPrice uint64, currency string) *Item {
	return &Item{
		name:          name,
		startingPrice: startingPrice,
		currentPrice:  startingPrice,
		currency:      currency,
	}
}

// UpdatePrice will update the price of the item
// if the new price is lower than old price return an error
// if the auction is done return an error
func (i *Item) UpdatePrice(newPrice uint64) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.locked {
		return ErrLocked
	}
	if i.currentPrice > newPrice {
		return ErrLower
	}
	i.currentPrice = newPrice
	return nil
}

// GetPrice will return the current item price
func (i *Item) GetPrice() uint64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.currentPrice
}

func (i *Item) GetName() string {
	return i.name
}

func (i *Item) GetCurrency() string {
	return i.currency
}

// Lock will lock the struct
// The use case of this is when the auction on the Item ends
func (i *Item) Lock() {
	i.locked = true
}
