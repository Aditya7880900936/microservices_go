package order

import (
	"context"
	"time"
)

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

type OrderedProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    uint32  `json:"quantity"`
}


type orderService struct {
	repository Repository
}

func NewService(r Repository) service {
	return &orderService{r}
}

func (s orderService) PostOrder(ctx context.Context , accountID string, products []OrderedProduct)(*Order, error) {
	order := &Order{
		ID:         "123",
		CreatedAt:  time.Now(),
		TotalPrice: 123.45,
		AccountID:  accountID,
		Products:   products,
	}
	return order, nil
}

func (s orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return nil, nil
}	