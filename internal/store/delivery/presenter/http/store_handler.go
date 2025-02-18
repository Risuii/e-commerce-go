package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/authentication/delivery/dto"
	StoreDTO "e-commerce/internal/store/delivery/dto"
	StoreUsecase "e-commerce/internal/store/domain/usecase"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	CustomValidationPackage "e-commerce/pkg/custom_validation"
	LoggerPackage "e-commerce/pkg/logger"
	RequestPackage "e-commerce/pkg/request_information"
)

type StoreHandler interface {
	CreateStore(*gin.Context)
	UpdateStore(*gin.Context)
	GetStore(*gin.Context)
}

type StoreHandlerImpl struct {
	library            Library.Library
	customValidation   CustomValidationPackage.CustomValidation
	createStoreUsecase StoreUsecase.CreateStoreUsecase
	updateStoreUsecase StoreUsecase.UpdateStoreUsecase
	getStoreUsecase    StoreUsecase.GetStoreUsecase
}

func NewStoreHandler(
	library Library.Library,
	customValidation CustomValidationPackage.CustomValidation,
	createStoreUsecase StoreUsecase.CreateStoreUsecase,
	updateStoreUsecase StoreUsecase.UpdateStoreUsecase,
	getStoreUsecase StoreUsecase.GetStoreUsecase,
) StoreHandler {
	return &StoreHandlerImpl{
		library:            library,
		customValidation:   customValidation,
		createStoreUsecase: createStoreUsecase,
		updateStoreUsecase: updateStoreUsecase,
		getStoreUsecase:    getStoreUsecase,
	}
}

func (h *StoreHandlerImpl) CreateStore(c *gin.Context) {
	path := "StoreHandler:CreateStore"

	var response *gin.H

	// INIT PARAM
	var param StoreDTO.StoreParam

	// GET REQUEST
	requestInformation := RequestPackage.RequestInformation{}
	request := requestInformation.GetRequestInformation(c)

	// GET TRACEID
	traceID, exists := c.Get(Constants.TraceID)
	if !exists {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, h.library)
		response = &gin.H{
			Constants.Path:    path,
			Constants.Message: Constants.ErrEmptyTraceID.Error(),
		}
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Request:  request,
			Constants.Response: response,
		}).Debug()

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// GET CREDENTIAL
	credentialPayload, exists := c.Get(Constants.Credential)
	if !exists {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, h.library)
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    path,
			Constants.Message: Constants.ErrEmptyCredential.Error(),
		}
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Request:  request,
			Constants.Response: response,
		}).Debug()

		c.JSON(http.StatusBadRequest, response)
		return
	}
	credential := credentialPayload.(*AuthDTO.LogoutParam)

	// CHECK VALIDATION
	errValidationPayload := param.Validate(requestInformation, h.library, h.customValidation)
	if len(errValidationPayload) > 0 {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, Constants.NilString, h.library)
		err = err.(*CustomErrorPackage.CustomError).FromListMap(errValidationPayload)
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    path,
			Constants.Message: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
			Constants.Data:    errValidationPayload,
		}

		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	// LOGIC USECASE
	usecase := h.createStoreUsecase
	err := usecase.Index(&param, credential)
	if err != nil {
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Message: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}

		c.JSON(err.(*CustomErrorPackage.CustomError).GetCode(), response)
		c.Abort()
		return
	}

	// RESPONSE
	response = &gin.H{
		Constants.TraceID: traceID,
		Constants.Path:    requestInformation.Path,
		Constants.Message: Constants.MsgSuccessSaveRequest,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *StoreHandlerImpl) UpdateStore(c *gin.Context) {
	path := "StoreHandler:UpdateStore"

	var response *gin.H

	// INIT PARAM
	var param StoreDTO.StoreParam

	// GET REQUEST
	requestInformation := RequestPackage.RequestInformation{}
	request := requestInformation.GetRequestInformation(c)

	// GET TRACEID
	traceID, exists := c.Get(Constants.TraceID)
	if !exists {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, h.library)
		response = &gin.H{
			Constants.Path:    path,
			Constants.Message: Constants.ErrEmptyTraceID.Error(),
		}
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Request:  request,
			Constants.Response: response,
		}).Debug()

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// GET CREDENTIAL
	credentialPayload, exists := c.Get(Constants.Credential)
	if !exists {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, h.library)
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    path,
			Constants.Message: Constants.ErrEmptyCredential.Error(),
		}
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Request:  request,
			Constants.Response: response,
		}).Debug()

		c.JSON(http.StatusBadRequest, response)
		return
	}
	credential := credentialPayload.(*AuthDTO.LogoutParam)

	// CHECK VALIDATION
	errValidationPayload := param.Validate(requestInformation, h.library, h.customValidation)
	if len(errValidationPayload) > 0 {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, Constants.NilString, h.library)
		err = err.(*CustomErrorPackage.CustomError).FromListMap(errValidationPayload)
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    path,
			Constants.Message: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
			Constants.Data:    errValidationPayload,
		}

		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	// LOGIC USECASE
	usecase := h.updateStoreUsecase
	err := usecase.Index(&param, credential)
	if err != nil {
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Message: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}

		c.JSON(err.(*CustomErrorPackage.CustomError).GetCode(), response)
		c.Abort()
		return
	}

	// RESPONSE
	response = &gin.H{
		Constants.TraceID: traceID,
		Constants.Path:    requestInformation.Path,
		Constants.Message: Constants.MsgSuccessSaveRequest,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *StoreHandlerImpl) GetStore(c *gin.Context) {
	path := "StoreHandler:GetStore"

	var response *gin.H

	// GET REQUEST
	requestInformation := RequestPackage.RequestInformation{}
	requestInformation.GetRequestInformation(c)

	// GET TRACEID
	traceID, exists := c.Get(Constants.TraceID)
	if !exists {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, h.library)
		response = &gin.H{
			Constants.Path:    path,
			Constants.Message: Constants.ErrEmptyTraceID.Error(),
		}
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Response: response,
		}).Debug()

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// GET CREDENTIAL
	credentialPayload, exists := c.Get(Constants.Credential)
	if !exists {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, h.library)
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    path,
			Constants.Message: Constants.ErrEmptyCredential.Error(),
		}
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Response: response,
		}).Debug()

		c.JSON(http.StatusBadRequest, response)
		return
	}
	credential := credentialPayload.(*AuthDTO.LogoutParam)

	// LOGIC USECASE
	usecase := h.getStoreUsecase
	responseDTO, err := usecase.Index(credential)
	if err != nil {
		response = &gin.H{
			Constants.TraceID: traceID,
			Constants.Path:    err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Message: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}

		c.JSON(err.(*CustomErrorPackage.CustomError).GetCode(), response)
		c.Abort()
		return
	}

	response = &gin.H{
		Constants.TraceID: traceID,
		Constants.Path:    requestInformation.Path,
		Constants.Message: Constants.MsgSuccessRequest,
		Constants.Data:    responseDTO,
	}

	c.JSON(http.StatusCreated, response)
}
