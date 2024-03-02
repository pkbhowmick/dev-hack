package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	sqlcdb "github.com/pkbhowmick/dev-hack/product/infra/sqlc/postgresql"
	productper "github.com/pkbhowmick/dev-hack/product/persistence/product"
	"github.com/pkbhowmick/dev-hack/product/usecase/product"
	"go.uber.org/dig"
)

func newDIContainer() (*dig.Container, error) {
	c := dig.New()

	pp := []interface{}{
		newSQLC,
		productper.NewRepository,
		newProductUsecase,
	}
	for _, p := range pp {
		if err := c.Provide(p); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func newProductUsecase(ur *productper.Repository) *product.Usecase {
	return &product.Usecase{
		ProductRepo: ur,
	}
}

func newSQLC() (*sqlcdb.Queries, error) {
	dburl := os.Getenv("POSTGRESQL_DATABASE_URL")
	if dburl == "" {
		return nil, errors.New("postgresql database url env is not set")
	}

	db, err := sql.Open("postgres", dburl)
	if err != nil {
		return nil, errors.New("couldn't open the postgres DB, because:" + err.Error())
	}

	var counter int
	for {
		if counter == 30 {
			log.Fatal("reached maximum number of attempt connecting to database")
		}

		fmt.Println("attempt to connect to postgres database", "counter", counter)
		err := db.Ping()
		if err == nil {
			break
		}

		log.Println("attempt connecting to postgres database failed, will be repeated in one second", "err", err)
		time.Sleep(time.Second)
		counter++
	}

	queries := sqlcdb.New(db)

	return queries, nil
}
