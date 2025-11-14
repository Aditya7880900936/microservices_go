package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO orders (id, created_at, total_price, account_id) VALUES ($1, $2, $3, $4)",
		o.ID, o.CreatedAt, o.TotalPrice, o.AccountID)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	if err != nil {
		return err
	}

	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	err = stmt.Close()
	return err
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {

	rows, err := r.db.QueryContext(ctx,
		`SELECT 
            o.id,
            o.created_at,
            o.account_id,
            o.total_price::money::numeric::float8,
            op.product_id,
            op.quantity
        FROM orders o
        JOIN order_products op ON o.id = op.order_id
        WHERE o.account_id = $1
        ORDER BY o.id`, accountID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	var currentOrder Order
	var products []OrderedProduct
	var lastOrderID string

	for rows.Next() {
		var (
			orderID    string
			createdAt  time.Time
			accID      string
			totalPrice float64
			productID  string
			quantity   int
		)

		err = rows.Scan(&orderID, &createdAt, &accID, &totalPrice, &productID, &quantity)
		if err != nil {
			return nil, err
		}

		// New order detected — append previous one
		if lastOrderID != "" && lastOrderID != orderID {
			currentOrder.Products = products
			orders = append(orders, currentOrder)
			products = []OrderedProduct{}
		}

		// Set new order header
		if lastOrderID != orderID {
			currentOrder = Order{
				ID:         orderID,
				CreatedAt:  createdAt,
				AccountID:  accID,
				TotalPrice: totalPrice,
			}
			lastOrderID = orderID
		}

		products = append(products, OrderedProduct{
			ID: productID,
			Quantity:  uint32(quantity),
		})
	}

	// Append last order if rows existed
	if lastOrderID != "" {
		currentOrder.Products = products
		orders = append(orders, currentOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
