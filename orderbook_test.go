package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l)

}

func TestOrderBook(t *testing.T) {
	ob := NewOrderbook()

	buyOrder := NewOrder(true, 10)
	ob.PlaceOrder(95_000, buyOrder)

	fmt.Println("Bids", ob.Bids)
}
