package account

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context , skip uint64 , take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err!= nil {
		db.Close()
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts (id, name) VALUES ($1, $2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	var a Account
	err := r.db.QueryRowContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id).Scan(&a.ID, &a.Name)
	if err!= nil {
		return nil, err
	}
	return &a, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows , err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",  skip, take)
	if err!= nil {
		return nil, err
	}
	defer rows.Close()
	account := []Account{}
	for rows.Next() {
		var a Account
		if err:= rows.Scan(&a.ID, &a.Name); err!= nil {
			return nil, err
		}
		account = append(account, a)
	}
	if err:= rows.Err(); err!= nil {
		return nil, err
	}
	return account, nil
}
