package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
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

var redisClient *redis.Client
var ctx = context.Background()

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (r *Repository) GetAllProductsByUserId(ctx context.Context, userId string) ([]sqlcdb.Product, error) {
	initRedis()
	val, err := redisClient.Get(ctx, userId).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err != redis.Nil {
		cachedData, err := convertStringToData(val)
		return cachedData, err
	}

	products, err := r.DB.GetProductsByUserId(ctx, userId)
	if err == sql.ErrNoRows {
		return nil, errors.New("user doesn't exist")
	}
	if err != nil {
		return nil, err
	}

	dataString, err := convertDataToString(products)
	if err != nil {
		return nil, err
	}

	// Set the data in Redis with an expiration (e.g., 1 hour)
	err = redisClient.Set(ctx, userId, dataString, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func convertStringToData(val string) ([]sqlcdb.Product, error) {
	var products []sqlcdb.Product
	err := json.Unmarshal([]byte(val), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func convertDataToString(products []sqlcdb.Product) (string, error) {
	jsonData, err := json.Marshal(products)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
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
