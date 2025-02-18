package dto

import (
	Library "e-commerce/library"
	CustomValidation "e-commerce/pkg/custom_validation"
	RequestPackage "e-commerce/pkg/request_information"
)

type StoreParam struct {
	StoreName   string `json:"store_name" validate:"required,alphanumeric"`
	Description string `json:"description" validate:"omitempty"`
}

func (e *StoreParam) Validate(request RequestPackage.RequestInformation, l Library.Library, customValidation CustomValidation.CustomValidation) []map[string]interface{} {
	validateStruct := customValidation.ConvertStructToInterfaceFields(e)
	if err := l.JsonUnmarshal([]byte(request.RequestBody), &validateStruct); err != nil {
		e = &StoreParam{}
	}
	errValidations := customValidation.ValidateStruct(validateStruct, "name")
	if len(errValidations) > 0 {
		return errValidations
	}

	err := l.JsonUnmarshal([]byte(request.RequestBody), &e)
	if err != nil {
		e = &StoreParam{}
	}

	return customValidation.ValidateStruct(e, "name")
}
