package product

import (
	"context"

	"github.com/pkbhowmick/dev-hack/product/domain/product"
	sqlcdb "github.com/pkbhowmick/dev-hack/product/infra/sqlc/postgresql"
)

type ProductRepository interface {
	GetAllProductsByUserId(ctx context.Context, userId string) ([]sqlcdb.Product, error)
	Create(ctx context.Context, opts *product.CreationOptions) (string, error)
}

type Usecase struct {
	ProductRepo ProductRepository
}

func (uc *Usecase) GetAllProducts(ctx context.Context, userId string) ([]sqlcdb.Product, error) {
	return uc.ProductRepo.GetAllProductsByUserId(ctx, userId)
}

func (uc *Usecase) CreateProduct(ctx context.Context, opts *product.CreationOptions) (string, error) {
	return uc.ProductRepo.Create(ctx, opts)
}
