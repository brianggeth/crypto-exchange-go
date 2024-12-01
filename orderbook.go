package main

import (
	"fmt"
	"time"
)

// Order is the struct of a Bid or Ask placed by the user at certain Limit (price).
type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

// Limit represents a certain price position, where different Orders are placed at.
type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: Orders{},
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}

	o.Limit = nil
	l.TotalVolume -= o.Size

	// TODO: resort the whole resting orders
}

func (l *Limit) String() string {
	return fmt.Sprintf("[price: %.2f] | volume: %.2f]", l.Price, l.TotalVolume)
}

// Orderbook represents the summatory of different possibles Limits (price area of interest) where Orders are placed at.
type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks: []*Limit{},
		Bids: []*Limit{},

		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

// Match is the struct of a match between a Bid and an Ask order placed at the same price.
// A match can be partial, meaning that only a part of the order is matched.
type Match struct {
	Ask   *Order
	Bid   *Order
	Size  float64
	Price float64
}

func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
	// 1. Try to match the orders (partial)
	// TODO: Matching logic

	// 2. Add the rest of the order to the books
	if o.Size > 0.0 {
		// TODO: All the rest of the order to the book
		ob.Add(price, o)
	}

	return []Match{}
}

func (ob *Orderbook) Add(price float64, o *Order) {
	var limit *Limit
	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	}

	limit.AddOrder(o)
}
