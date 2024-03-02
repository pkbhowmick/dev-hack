package signup

import (
	"context"

	productdom "github.com/pkbhowmick/dev-hack/gateway/domain/product"
	sqlcdb "github.com/pkbhowmick/dev-hack/gateway/infra/sqlc/mysql"
)

type ProductRepository interface {
	GetUserById(ctx context.Context, userId string) (*sqlcdb.User, error)
}

type Usecase struct {
	ProductRepository ProductRepository
}

func (uc *Usecase) CreateProduct(ctx context.Context, opts *productdom.CreationOptions, userId string) (string, error) {
	// call the product service to create the product

	return "", nil // return productId that get from the product service call
}

func (uc *Usecase) GetAllProducts(ctx context.Context, userId string) ([]productdom.Product, error) {
	// call the product service to get the products by the userid

	return nil, nil // return all products of that user that get from the product service call
}
