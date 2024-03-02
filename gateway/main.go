package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"github.com/pkbhowmick/dev-hack/gateway/presenter/healthcheck"
	"github.com/pkbhowmick/dev-hack/gateway/presenter/product"
	"github.com/pkbhowmick/dev-hack/gateway/presenter/signup"
	productuc "github.com/pkbhowmick/dev-hack/gateway/usecase/product"
	signupuc "github.com/pkbhowmick/dev-hack/gateway/usecase/signup"
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

	log.Println("server is running")
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
		signupuc *signupuc.Usecase,
		productuc *productuc.Usecase,
	) {
		r.GET("/", healthcheck.Handler())
		r.POST("/signup", signup.Handler(signupuc))
		r.POST("/users/:userId/products", product.CreationHandler(productuc))
		r.GET("/users/:userId/products", product.ListHandler(productuc))
	})
	if err != nil {
		return nil, err
	}

	return &http.Server{Addr: ":8080", Handler: r}, nil
}
