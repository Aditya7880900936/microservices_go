package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
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

func (s orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	o := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}
	o.TotalPrice = 0
	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := s.repository.PutOrder(ctx, *o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	orders, err := s.repository.GetOrdersForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
