package dto

import (
	Library "e-commerce/library"
	CustomValidation "e-commerce/pkg/custom_validation"
	RequestPackage "e-commerce/pkg/request_information"
)

type RegisterParam struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanumeric,max=15"`
	Password string `json:"password" validate:"required"`
}

func (e *RegisterParam) Validate(request RequestPackage.RequestInformation, l Library.Library, customValidation CustomValidation.CustomValidation) []map[string]interface{} {
	validateStruct := customValidation.ConvertStructToInterfaceFields(e)
	if err := l.JsonUnmarshal([]byte(request.RequestBody), &validateStruct); err != nil {
		e = &RegisterParam{}
	}
	errValidations := customValidation.ValidateStruct(validateStruct, "name")
	if len(errValidations) > 0 {
		return errValidations
	}

	err := l.JsonUnmarshal([]byte(request.RequestBody), &e)
	if err != nil {
		e = &RegisterParam{}
	}

	return customValidation.ValidateStruct(e, "name")
}
