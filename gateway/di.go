package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	sqlcdb "github.com/pkbhowmick/dev-hack/gateway/infra/sqlc/mysql"
	userpersistence "github.com/pkbhowmick/dev-hack/gateway/persistence/user"
	"github.com/pkbhowmick/dev-hack/gateway/usecase/signup"
	"go.uber.org/dig"
)

func newDIContainer() (*dig.Container, error) {
	c := dig.New()

	pp := []interface{}{
		newSQLC,
		newUserRepository,
		newSignupUsecase,
	}
	for _, p := range pp {
		if err := c.Provide(p); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func newUserRepository(db *sqlcdb.Queries) *userpersistence.Repository {
	return &userpersistence.Repository{DB: db}
}

func newSignupUsecase(ur *userpersistence.Repository) *signup.Usecase {
	return &signup.Usecase{
		UserRepository: ur,
	}
}

func newSQLC() (*sqlcdb.Queries, error) {
	dburl := os.Getenv("MYSQL_DATABASE_URL")
	if dburl == "" {
		return nil, errors.New("mysql database url env is not set")
	}

	db, err := sql.Open("postgres", dburl)
	if err != nil {
		return nil, errors.New("couldn't open the mysql DB, because:" + err.Error())
	}

	var counter int
	for {
		if counter == 30 {
			log.Fatal("reached maximum number of attempt connecting to database")
		}

		fmt.Println("attempt to connect to mysql database", "counter", counter)
		err := db.Ping()
		if err == nil {
			break
		}

		log.Println("attempt connecting to mysql database failed, will be repeated in one second", "err", err)
		time.Sleep(time.Second)
		counter++
	}

	queries := sqlcdb.New(db)

	return queries, nil
}
