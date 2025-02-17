package usecase

import (
	"time"

	Constants "e-commerce/constants"
	LogDTO "e-commerce/internal/logging/delivery/dto"
	LogEntity "e-commerce/internal/logging/domain/entity"
	LogRepository "e-commerce/internal/logging/domain/repository"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	UtilsPackage "e-commerce/pkg/utils"
)

type LoggingUsecase interface {
	Index(traceID, pathURL string, request, response []byte) (*LogDTO.LogActivityParam, error)
}

type LoggingUsecaseImpl struct {
	logRepository LogRepository.LogActivity
	library       Library.Library
}

func NewLogUsecase(
	logRepository LogRepository.LogActivity,
	library Library.Library,
) LoggingUsecase {
	return &LoggingUsecaseImpl{
		logRepository: logRepository,
		library:       library,
	}
}

func (u *LoggingUsecaseImpl) Index(traceID, pathURL string, request, response []byte) (*LogDTO.LogActivityParam, error) {
	path := "LogActivityUsecase:Index"

	// UNMARSHAL RESPONSE
	var logData LogDTO.LogActivityParam
	err := u.library.JsonUnmarshal(response, &logData)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// MARSHAL FIELD DATA RESPONSE
	responseData, err := u.library.JsonMarshal(logData.Data)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// REMOVE REQUEST IF RESPONSE FROM API REGISTER OR LOGIN
	if pathURL == "/api/v1/user/register" || pathURL == "/api/v1/user/login" {
		request = nil
	}

	// SET ENTITY LOGGING
	dataLog := LogEntity.LogActivity{
		TraceID:         traceID,
		Endpoint:        pathURL,
		Path:            logData.Path,
		Description:     logData.Message,
		CreatedAt:       time.Now().Format(Constants.YYYYMMDDHHmmss),
		RequestPayload:  string(request),
		ResponsePayload: UtilsPackage.TernaryOperator(logData.Data == nil, string(response), string(responseData)),
	}

	// INSERT ENTITY LOGGING TO DB
	go u.logRepository.CreateLog(dataLog)

	return &logData, nil
}
