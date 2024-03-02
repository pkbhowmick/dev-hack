package signup

import (
	"context"

	productdom "github.com/pkbhowmick/dev-hack/gateway/domain/product"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) error
}

type Usecase struct {
	UserRepository UserRepository
}

func (uc *Usecase) CreateProduct(ctx context.Context, opts *productdom.CreationOptions, userId string) (string, error) {
	if err := uc.UserRepository.GetUserById(ctx, userId); err != nil {
		return "", err
	}

	// call the product service to create the product

	return "", nil // return productId that get from the product service call
}

func (uc *Usecase) GetAllProducts(ctx context.Context, userId string) ([]productdom.Product, error) {
	if err := uc.UserRepository.GetUserById(ctx, userId); err != nil {
		return nil, err
	}

	// call the product service to get the products by the userid

	return nil, nil // return all products of that user that get from the product service call
}
