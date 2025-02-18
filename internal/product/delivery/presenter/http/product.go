package http

import (
	"github.com/gin-gonic/gin"

	Library "e-commerce/library"
	CustomValidationPackage "e-commerce/pkg/custom_validation"
)

type ProductHandler interface {
	CreateProduct(*gin.Context)
}

type ProductHandlerImpl struct {
	library          Library.Library
	customValidation CustomValidationPackage.CustomValidation
}

func NewProductHandler(
	library Library.Library,
	customValidation CustomValidationPackage.CustomValidation,
) ProductHandler {
	return &ProductHandlerImpl{
		library:          library,
		customValidation: customValidation,
	}
}

func (h *ProductHandlerImpl) CreateProduct(c *gin.Context) {}
