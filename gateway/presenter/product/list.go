package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/is"
	. "github.com/go-ozzo/ozzo-validation/v4"
	productdom "github.com/pkbhowmick/dev-hack/gateway/domain/product"
	productuc "github.com/pkbhowmick/dev-hack/gateway/usecase/product"
)

type response struct {
	OK       bool                 `json:"ok"`
	Message  string               `json:"message"`
	Products []productdom.Product `json:"products"`
}

func ListHandler(uc *productuc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		err := Validate(userId, Required, is.UUID)
		if err != nil {
			returnCreationResponse(c, http.StatusBadRequest, false, err.Error(), "")
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

func returnResponse(c *gin.Context, statusCode int, ok bool, message string, products []productdom.Product) {
	c.JSON(statusCode, &response{
		OK:       ok,
		Message:  message,
		Products: products,
	})
}
