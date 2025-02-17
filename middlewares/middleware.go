package middlewares

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	Constants "e-commerce/constants"
	LoggingUsecase "e-commerce/internal/logging/domain/usecase"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	LoggerPackage "e-commerce/pkg/logger"
	RequestPackage "e-commerce/pkg/request_information"
	ResponsePackage "e-commerce/pkg/response_information"
)

type Middleware interface {
	GenerateTraceID() gin.HandlerFunc
	Logging() gin.HandlerFunc
}

type MiddlewareImpl struct {
	library        Library.Library
	loggingUsecase LoggingUsecase.LoggingUsecase
}

func NewMiddleware(
	library Library.Library,
	loggingUsecase LoggingUsecase.LoggingUsecase,
) Middleware {
	return &MiddlewareImpl{
		library:        library,
		loggingUsecase: loggingUsecase,
	}
}

func (m *MiddlewareImpl) GenerateTraceID() gin.HandlerFunc {
	path := "Middleware:GenerateTraceID"
	return func(c *gin.Context) {
		var response *gin.H
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

		c.Set(Constants.TraceID, traceID)

		c.Next()
	}
}

func (m *MiddlewareImpl) Logging() gin.HandlerFunc {
	path := "Middleware:Logging"

	return func(c *gin.Context) {
		buffer := &bytes.Buffer{}
		writer := &ResponsePackage.ResponseWriter{Body: buffer, ResponseWriter: c.Writer}
		c.Writer = writer

		requestInformation := RequestPackage.RequestInformation{}
		request := requestInformation.GetRequestInformation(c)

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

		c.Next()

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
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(writer.StatusCode)
		writer.WriteResponse([]byte(response))
	}
}
