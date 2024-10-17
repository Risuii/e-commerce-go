package logger

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	Constants "e-commerce/constants"
	LogEntity "e-commerce/internal/logging/domain/entity"
	Library "e-commerce/library"

	Logging "e-commerce/internal/logging/domain/repository"
)

var logger *logrus.Logger

func New(
	library Library.Library,
) {
	logLevel := logrus.DebugLevel
	log := logrus.New()
	log.SetLevel(logLevel)

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		PrettyPrint:       true,
		DisableHTMLEscape: true,
	})

	logger = log
}

type Logger interface {
	LoggingActivty(logData LogEntity.LogActivity)
	LoggingInterface(logData LogEntity.LogInterface)
	LoggingOutgoing(logData LogEntity.LogOutgoing)
}

type LoggerImpl struct {
	logActivity  Logging.LogActivity
	logInterface Logging.LogInterface
	logOutgoing  Logging.LogOutgoing
}

func NewLogger(logActivity Logging.LogActivity, logInterface Logging.LogInterface, logOutgoing Logging.LogOutgoing) Logger {
	return &LoggerImpl{
		logActivity:  logActivity,
		logInterface: logInterface,
		logOutgoing:  logOutgoing,
	}
}

func (l LoggerImpl) LoggingActivty(logData LogEntity.LogActivity) {
	var traceId string

	date := time.Now().Format(Constants.YYYMMDDHHMMSS)

	if logData.TraceId == Constants.NilString {
		traceId, _ = Library.New().GenerateUUID()
	} else {
		traceId = logData.TraceId
	}

	data := LogEntity.LogActivity{
		TraceId:        traceId,
		Endpoint:       logData.Endpoint,
		Path:           logData.Path,
		Description:    logData.Description,
		CreatedAt:      date,
		RequestPayload: logData.RequestPayload,
	}

	fmt.Println(time.Now().Format(Constants.YYYMMDDHHMMSS), Constants.TraceIdActivity, traceId)

	l.logActivity.CreateLog(data)
}

func (l LoggerImpl) LoggingInterface(logData LogEntity.LogInterface) {
	date := time.Now().Format(Constants.YYYMMDDHHMMSS)

	data := LogEntity.LogInterface{
		TraceId:         logData.TraceId,
		ServerNode:      logData.ServerNode,
		ServiceName:     logData.ServiceName,
		RequestPayload:  logData.RequestPayload,
		RequestDate:     logData.RequestDate,
		ResponsePayload: logData.ResponsePayload,
		ResponseDate:    date,
		SourceName:      logData.SourceName,
		MsName:          logData.MsName,
	}

	fmt.Println(time.Now().Format(Constants.YYYMMDDHHMMSS), Constants.TraceIdInterface, logData.TraceId)

	l.logInterface.CreateLog(data)

}

func (l LoggerImpl) LoggingOutgoing(logData LogEntity.LogOutgoing) {
	// date := time.Now().Format(Constants.YYYMMDDHHMMSS)

	// data := LogEntity.LogOutgoing{
	// 	TraceId:         logData.TraceId,
	// 	BackendSystem:   logData.BackendSystem,
	// 	ServiceName:     logData.ServiceName,
	// 	RequestPayload:  logData.RequestPayload,
	// 	ResponsePayload: logData.ResponsePayload,
	// 	RequestDate:     logData.RequestDate,
	// 	ResponseDate:    date,
	// 	StatusCode:      logData.StatusCode,
	// 	MsName:          logData.MsName,
	// 	ServerNode:      logData.ServerNode,
	// 	SourceNode:      logData.SourceNode,
	// 	ResponseTime:    logData.ResponseTime,
	// }

	fmt.Println(time.Now().Format(Constants.YYYMMDDHHMMSS), Constants.TraceIdOutgoing, logData.TraceId)

	// l.logOutgoing.CreateLog(data)
}

func WriteLog(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}
