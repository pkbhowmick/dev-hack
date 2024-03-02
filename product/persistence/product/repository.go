package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/pkbhowmick/dev-hack/product/domain/product"
	sqlcdb "github.com/pkbhowmick/dev-hack/product/infra/sqlc/postgresql"
)

type Repository struct {
	DB *sqlcdb.Queries
}

func NewRepository(db *sqlcdb.Queries) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetAllProductsByUserId(ctx context.Context, userId string) ([]sqlcdb.Product, error) {
	products, err := r.DB.GetProductsByUserId(ctx, userId)
	if err == sql.ErrNoRows {
		return nil, errors.New("user doesn't exist")
	}
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) Create(ctx context.Context, opts *product.CreationOptions) (string, error) {
	arg := sqlcdb.CreateProductParams{
		ID:          uuid.NewString(),
		UserID:      opts.UserId,
		Name:        opts.Name,
		Price:       int32(opts.Price),
		Description: opts.Description,
		Stock:       int32(opts.Stock),
	}
	err := r.DB.CreateProduct(ctx, arg)
	if err != nil {
		return "", err
	}

	return arg.ID, nil
}
