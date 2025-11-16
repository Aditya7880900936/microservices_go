package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    if id != nil {
        acc, err := r.server.accountClient.GetAccount(ctx, *id)
        if err != nil {
            log.Println(err)
            return nil, err
        }
        return []*Account{
            {
                ID:   acc.ID,
                Name: acc.Name,
            },
        }, nil
    }

    skip, take := uint64(0), uint64(0)
    if pagination != nil {
        skip, take = pagination.bounds()
    }

    accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // FIXED: create correct slice
    var accounts []*Account

    for _, acc := range accountList {
        accounts = append(accounts, &Account{
            ID:   acc.ID,
            Name: acc.Name,
        })
    }

    return accounts, nil
}


func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, canel := context.WithTimeout(ctx, 3*time.Second)
	defer canel()
	if id!= nil {
		prod, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err!= nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{
			{
				ID:   prod.ID,
				Name: prod.Name,
				Price: prod.Price,
				Description: prod.Description,
			}}, nil
	}

	skip, take := uint64(0),  uint64(0)

	if pagination!= nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query!= nil {
		q = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, skip , take , nil , q)
	if err!= nil {
		log.Println(err)
		return nil, err
	}
	var products []*Product
	for _, prod := range productList{
		products = append(products , &Product{
			ID:   prod.ID,
			Name: prod.Name,
			Price: prod.Price,
			Description: prod.Description,
		},
	)
	}
	return products, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(0)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Limit != nil {
		takeValue = uint64(*p.Limit)
	}
	return skipValue, takeValue
}
