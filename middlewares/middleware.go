package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/authentication/delivery/dto"
	LoggingUsecase "e-commerce/internal/logging/domain/usecase"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	JWEPackage "e-commerce/pkg/jwe"
	LoggerPackage "e-commerce/pkg/logger"
	RequestPackage "e-commerce/pkg/request_information"
	ResponsePackage "e-commerce/pkg/response_information"
)

type Middleware interface {
	GenerateTraceID() gin.HandlerFunc
	Logging() gin.HandlerFunc
	ValidateToken() gin.HandlerFunc
}

type MiddlewareImpl struct {
	config         Config.Config
	library        Library.Library
	jwe            JWEPackage.JWE
	loggingUsecase LoggingUsecase.LoggingUsecase
}

func NewMiddleware(
	config Config.Config,
	library Library.Library,
	jwe JWEPackage.JWE,
	loggingUsecase LoggingUsecase.LoggingUsecase,
) Middleware {
	return &MiddlewareImpl{
		config:         config,
		library:        library,
		jwe:            jwe,
		loggingUsecase: loggingUsecase,
	}
}

func (m *MiddlewareImpl) GenerateTraceID() gin.HandlerFunc {
	path := "Middleware:GenerateTraceID"
	return func(c *gin.Context) {
		var response *gin.H

		// GENERATE UUID FOR TRACEID
		traceID, err := m.library.GenerateUUID()
		if err != nil {
			err := CustomErrorPackage.New(Constants.ErrFailedGenerateTraceID, Constants.ErrFailedGenerateTraceID, path, m.library)

			response = &gin.H{
				Constants.Status:     true,
				Constants.CodeStatus: Constants.StatusAuthenticationFailuer,
				Constants.Message:    Constants.ErrUnauthorized.Error(),
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Response: response,
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// SET TRACEID TO GIN CONTEXT
		c.Set(Constants.TraceID, traceID)

		c.Next()
	}
}

func (m *MiddlewareImpl) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := "Middleware:ValidateToken"

		var response *gin.H

		requestInformation := RequestPackage.RequestInformation{}
		request := requestInformation.GetRequestInformation(c)

		// GET TRACEID
		traceID, exists := c.Get(Constants.TraceID)
		if !exists {
			err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, m.library)
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

		// GET TOKEN FROM HEADER
		authenticationHeader := c.GetHeader(Constants.Authorization)

		// CHECK TOKEN IS EXIST
		if authenticationHeader == Constants.NilString {
			err := CustomErrorPackage.New(Constants.ErrInvalidJWE, Constants.ErrInvalidJWE, path, m.library)

			response = &gin.H{
				Constants.TraceID: traceID,
				Constants.Message: Constants.ErrUnauthorized.Error(),
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Request:  request,
				Constants.Response: response,
			})
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// CHECKING PREFIX TOKEN
		if !strings.HasPrefix(authenticationHeader, Constants.Bearer) {
			err := CustomErrorPackage.New(Constants.ErrInvalidJWE, Constants.ErrInvalidJWE, path, m.library)

			response = &gin.H{
				Constants.TraceID: traceID,
				Constants.Message: Constants.ErrAuthorizationBearer.Error(),
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Request:  request,
				Constants.Response: response,
			})
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// REMOVE PREFIX
		token := strings.TrimPrefix(authenticationHeader, Constants.Bearer)

		// VALIDATE TOKEN
		credential, err := m.jwe.JWEValidateToken([]byte(token), m.config.GetConfig().JWE.SecretKey)
		if string(credential) == Constants.NilString || err != nil {
			err := CustomErrorPackage.New(Constants.ErrInvalidJWE, Constants.ErrInvalidJWE, path, m.library)

			response = &gin.H{
				Constants.TraceID: traceID,
				Constants.Message: Constants.ErrUnauthorized.Error(),
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Request:  request,
				Constants.Response: response,
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		var credentials *AuthDTO.LogoutParam
		if err := m.library.JsonUnmarshal(credential, &credentials); err != nil {
			err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, m.library)
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

		// SET CREDENTIAL TO CONTEXT
		c.Set(Constants.Credential, credentials)

		// CONTINUE TO API
		c.Next()
	}
}

func (m *MiddlewareImpl) Logging() gin.HandlerFunc {
	path := "Middleware:Logging"

	return func(c *gin.Context) {
		// CATCH RESPONSE
		buffer := &bytes.Buffer{}
		writer := &ResponsePackage.ResponseWriter{Body: buffer, ResponseWriter: c.Writer}
		c.Writer = writer

		// GET REQUEST
		requestInformation := RequestPackage.RequestInformation{}
		request := requestInformation.GetRequestInformation(c)

		// GET TRACEID FROM CONTEXT
		traceID, exists := c.Get(Constants.TraceID)
		if !exists {
			err := CustomErrorPackage.New(Constants.ErrValidation, nil, path, m.library)
			response := gin.H{
				Constants.TraceID: traceID,
				Constants.Status:  false,
				Constants.Message: Constants.ErrEmptyTraceID.Error(),
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Request:  request,
				Constants.Response: response,
			})

			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// DO THE API
		c.Next()

		// GET RESPONSE FROM API
		plainResponse := buffer.Bytes()
		if len(plainResponse) == 0 {
			err := CustomErrorPackage.New(Constants.ErrInternalServerError, nil, path, m.library)
			response := gin.H{
				Constants.TraceID: traceID,
				Constants.Status:  false,
				Constants.Message: Constants.ErrFailedEncryptResponse,
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Request:  request,
				Constants.Response: response,
			})

			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		// DO LOGIC LOGGING
		responseData, err := m.loggingUsecase.Index(traceID.(string), requestInformation.Path, requestInformation.RequestBody, plainResponse)
		if err != nil {
			err := CustomErrorPackage.New(Constants.ErrInternalServerError, err, path, m.library)
			response := gin.H{
				Constants.TraceID: traceID,
				Constants.Status:  false,
				Constants.Message: Constants.ErrFailedEncryptResponse,
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Request:  request,
				Constants.Response: response,
			})

			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		// REMOVE FIELD PATH FROM RESPONSE
		responseData.Path = Constants.NilString
		response, err := m.library.JsonMarshal(responseData)
		if err != nil {
			err := CustomErrorPackage.New(Constants.ErrInternalServerError, err, path, m.library)
			response := gin.H{
				Constants.TraceID: traceID,
				Constants.Status:  false,
				Constants.Message: Constants.ErrFailedEncryptResponse,
			}

			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:     err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Request:  request,
				Constants.Response: response,
			})

			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		// WRITE NEW RESPONSE
		c.Writer.Header().Set(Constants.ContentType, Constants.ApplicationJson)
		c.Writer.WriteHeader(writer.StatusCode)
		fmt.Println(time.Now().Format(Constants.YYYMMDDHHMMSS), Constants.TraceID, traceID)
		writer.WriteResponse([]byte(response))
	}
}
