package example1

import (
	"fmt"

	"github.com/orzkratos/errkratos"
)

type Guest struct {
	Name    string
	GuestID int `json:"-"`
}

type GuestOrdersStates struct {
	Guest       *Guest
	OrderStates []*OrderState
	Outline     string
	Erk         *errkratos.Erk
}

type Order struct {
	Name    string
	OrderID int `json:"-"`
}

type OrderState struct {
	Order *Order
	State string
	Erk   *errkratos.Erk
}

func NewGuests(guestCount int) []*Guest {
	var guests = make([]*Guest, 0, guestCount)
	for idx := 0; idx < guestCount; idx++ {
		guests = append(guests, &Guest{
			Name:    fmt.Sprintf("guest(%d)", idx),
			GuestID: idx,
		})
	}
	return guests
}

func NewOrders(guest *Guest, orderCount int) []*Order {
	orders := make([]*Order, 0, orderCount)
	for idx := 0; idx < orderCount; idx++ {
		orders = append(orders, &Order{
			Name:    fmt.Sprintf("guest(%d) order(%d)", guest.GuestID, idx),
			OrderID: idx,
		})
	}
	return orders
}
