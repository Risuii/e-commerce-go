package dto

import (
	Library "e-commerce/library"
	CustomValidation "e-commerce/pkg/custom_validation"
	RequestPackage "e-commerce/pkg/request_information"
)

type StoreStatusParam struct {
	Status string `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
}

func (e *StoreStatusParam) Validate(request RequestPackage.RequestInformation, l Library.Library, customValidation CustomValidation.CustomValidation) []map[string]interface{} {
	validateStruct := customValidation.ConvertStructToInterfaceFields(e)
	if err := l.JsonUnmarshal([]byte(request.RequestBody), &validateStruct); err != nil {
		e = &StoreStatusParam{}
	}
	errValidations := customValidation.ValidateStruct(validateStruct, "name")
	if len(errValidations) > 0 {
		return errValidations
	}

	err := l.JsonUnmarshal([]byte(request.RequestBody), &e)
	if err != nil {
		e = &StoreStatusParam{}
	}

	return customValidation.ValidateStruct(e, "name")
}
