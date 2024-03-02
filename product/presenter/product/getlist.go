package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/is"
	. "github.com/go-ozzo/ozzo-validation/v4"
	sqlcdb "github.com/pkbhowmick/dev-hack/product/infra/sqlc/postgresql"
	"github.com/pkbhowmick/dev-hack/product/usecase/product"
)

type response struct {
	OK       bool             `json:"ok"`
	Message  string           `json:"message"`
	Products []sqlcdb.Product `json:"products"`
}

func GetListHandler(uc *product.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		err := Validate(userId, Required, is.UUID)
		if err != nil {
			returnResponse(c, http.StatusBadRequest, false, err.Error(), nil)
			return
		}

		products, err := uc.GetAllProducts(c, userId)
		if err != nil {
			returnResponse(c, http.StatusBadRequest, false, err.Error(), nil)
			return
		}

		returnResponse(c, http.StatusOK, true, "Products fetching is succeeded", products)
	}
}

func returnResponse(c *gin.Context, statusCode int, ok bool, message string, products []sqlcdb.Product) {
	c.JSON(statusCode, &response{
		OK:       ok,
		Message:  message,
		Products: products,
	})
}
