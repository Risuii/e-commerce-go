package dto

import (
	Library "e-commerce/library"
	CustomValidation "e-commerce/pkg/custom_validation"
	RequestPackage "e-commerce/pkg/request_information"
)

type LoginParam struct {
	Email    string `json:"email" validate:"email_required_without=Username"`
	Username string `json:"username" validate:"alphanumeric_required_without=Email"`
	Password string `json:"password" validate:"required"`
}

func (e *LoginParam) Validate(request RequestPackage.RequestInformation, l Library.Library, customValidation CustomValidation.CustomValidation) []map[string]interface{} {
	validateStruct := customValidation.ConvertStructToInterfaceFields(e)
	if err := l.JsonUnmarshal([]byte(request.RequestBody), &validateStruct); err != nil {
		e = &LoginParam{}
	}
	errValidations := customValidation.ValidateStruct(validateStruct, "name")
	if len(errValidations) > 0 {
		return errValidations
	}

	err := l.JsonUnmarshal([]byte(request.RequestBody), &e)
	if err != nil {
		e = &LoginParam{}
	}

	return customValidation.ValidateStruct(e, "name")
}
