package item

import "sync"

type Item struct {
	mu            sync.Mutex
	Name          string
	StartingPrice uint64
	CurrentPrice  uint64
	Currency      string
}

func New(name string, startingPrice uint64, currency string) *Item {
	return &Item{
		Name:          name,
		StartingPrice: startingPrice,
		CurrentPrice:  startingPrice,
		Currency:      currency,
	}
}
