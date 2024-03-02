package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/pkbhowmick/dev-hack/product/presenter/healthcheck"
	"github.com/pkbhowmick/dev-hack/product/presenter/product"
	productuc "github.com/pkbhowmick/dev-hack/product/usecase/product"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	server, err := newServer()
	if err != nil {
		return err
	}

	log.Println("server is running...")
	return server.ListenAndServe()
}

func newServer() (*http.Server, error) {
	c, err := newDIContainer()
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	err = c.Invoke(func(
		productUC *productuc.Usecase,
	) {
		r.GET("/", healthcheck.Handler())

		r.GET("/products/:userId", product.GetListHandler(productUC))
		r.POST("/products", product.CreationHandler(productUC))
	})
	if err != nil {
		return nil, err
	}

	return &http.Server{Addr: ":8080", Handler: r}, nil
}
