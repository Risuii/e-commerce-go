package utils

import (
	"github.com/sirupsen/logrus"

	Constants "e-commerce/constants"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	LoggerPackage "e-commerce/pkg/logger"
)

func CatchPanic(path string, library Library.Library) {
	if err := recover(); err != nil {
		err = CustomErrorPackage.New(Constants.ErrPanic, err.(error), path, library)
		LoggerPackage.WriteLog(logrus.Fields{
			"path":  err.(*CustomErrorPackage.CustomError).GetPath(),
			"title": err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}).Panic(err.(*CustomErrorPackage.CustomError).GetPlain())
	}
}

func TernaryOperator[T interface{}](comparator bool, trueCondition T, falseCondition T) T {
	if comparator {
		return trueCondition
	}

	return falseCondition
}
