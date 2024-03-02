package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/is"
	. "github.com/go-ozzo/ozzo-validation/v4"
	productdom "github.com/pkbhowmick/dev-hack/gateway/domain/product"
	productuc "github.com/pkbhowmick/dev-hack/gateway/usecase/product"
)

type creationResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	ProductID string `json:"productID"`
}

func CreationHandler(uc *productuc.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		err := Validate(userId, Required, is.UUID)
		if err != nil {
			returnCreationResponse(c, http.StatusBadRequest, false, err.Error(), "")
			return
		}

		r := new(productdom.CreationOptions)
		if err := c.ShouldBindJSON(&r); err != nil {
			returnCreationResponse(c, http.StatusBadRequest, false, err.Error(), "")
			return
		}

		productID, err := uc.CreateProduct(c, r, userId)
		if err != nil {
			returnCreationResponse(c, http.StatusBadRequest, false, err.Error(), "")
			return
		}

		returnCreationResponse(c, http.StatusOK, true, "Product is created successfully", productID)
	}
}

func returnCreationResponse(c *gin.Context, statusCode int, ok bool, message, productID string) {
	c.JSON(statusCode, &creationResponse{
		OK:        ok,
		Message:   message,
		ProductID: productID,
	})
}
