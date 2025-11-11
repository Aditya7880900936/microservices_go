package catalog

import (
	"context"
	"errors"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("Entity Not Found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}


func NewElasticRepository(url string) (Repository, error) {
	client , err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {
	r.client.Stop()
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		Id(p.ID).
		BodyJson(productDocument{
			Name: p.Name,
			Description: p.Description,
			Price: p.Price,
		}).
		Do(ctx)

	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	
}