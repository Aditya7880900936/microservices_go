package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name string, description string, price string) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (s *catalogService) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	product := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String(),
	}
	if err := s.repository.PutProduct(ctx, *product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	product, err := s.repository.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	products, err := s.repository.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *catalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {
	products, err := s.repository.ListProductsWithIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	products, err := s.repository.SearchProducts(ctx, query, skip, take)
	if err != nil {
		return nil, err
	}
	return products, nil
}
