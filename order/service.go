package order

import "time"

type service interface {
	PutOrder
	GetOrdersForAccount
}

type Order struct {
	ID         string           `json:"id"`
	CreatedAt  time.Time        `json:"createdAt"`
	TotalPrice float64          `json:"totalPrice"`
	AccountID  string           `json:"accountId"`
	Products   []OrderedProduct `json:"products"`
}
