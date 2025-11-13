package catalog

import (
	"context"

	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.catalogServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		service: pb.NewcatalogServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	r, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Product.ID,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		ID: id,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Product.ID,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]*Product, error) {
	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	})
	if err!= nil {
		return nil, err
	}
	products := []Product{}
	for _, p := range r.Products {
		products = append(products, Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
}
