package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/auth/delivery/dto"
	UserUseCase "e-commerce/internal/auth/domain/usecase"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	CustomValidationPackage "e-commerce/pkg/custom_validation"
	LoggerPackage "e-commerce/pkg/logger"
	RequestPackage "e-commerce/pkg/request_information"
)

type UserHandler interface {
	Register(*gin.Context)
	Login(*gin.Context)
}

type UserHandlerImpl struct {
	library          Library.Library
	customValidation CustomValidationPackage.CustomValidation
	userUseCase      UserUseCase.RegisterUseCase
}

func NewUserHandler(library Library.Library,
	customValidation CustomValidationPackage.CustomValidation,
	userUseCase UserUseCase.RegisterUseCase,
) UserHandler {
	return &UserHandlerImpl{
		library:          library,
		customValidation: customValidation,
		userUseCase:      userUseCase,
	}
}

func (h *UserHandlerImpl) Register(c *gin.Context) {
	path := Constants.UserRegisterHandler

	var response *gin.H

	// INIT PARAM
	var param AuthDTO.RegisterParam

	// GET REQUEST
	requestInformation := RequestPackage.RequestInformation{}
	request := requestInformation.GetRequestInformation(c)

	// GET TRACEID
	traceID, exists := c.Get(Constants.TraceID)
	if !exists {
		err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, h.library)
		response = &gin.H{
			Constants.TraceID: traceID,
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

	// CHECK VALIDATION
	errValidationPayload := param.Validate(requestInformation, h.library, h.customValidation)
	if len(errValidationPayload) > 0 {
		param.Password = Constants.NilString
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
	usecase := h.userUseCase
	err := usecase.Index(requestInformation, param, traceID.(string))
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
		Constants.TraceID:    traceID,
		Constants.CodeStatus: Constants.RequestSuccess,
		Constants.Path:       requestInformation.Path,
		Constants.Message:    Constants.MsgSuccessSaveRequest,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandlerImpl) Login(c *gin.Context) {}
