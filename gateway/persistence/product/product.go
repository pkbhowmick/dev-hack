package product

import (
	"context"

	sqlcdb "github.com/pkbhowmick/dev-hack/gateway/infra/sqlc/mysql"
)

type Repository struct {
	DB *sqlcdb.Queries
}

func (r *Repository) GetByName(ctx context.Context, name string) (sqlcdb.User, error) {
	return r.DB.GetUserByName(ctx, name)
}

func (r *Repository) GetUserById(ctx context.Context, id string) (*sqlcdb.User, error) {
	user, err := r.DB.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
