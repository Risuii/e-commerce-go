package repository

import (
	LogInterfaceSource "e-commerce/internal/logging/data/source"
	LogParamEntity "e-commerce/internal/logging/domain/entity"
	LogInterfaceRepository "e-commerce/internal/logging/domain/repository"
)

type LogInterfaceImpl struct {
	logInterfaceSource LogInterfaceSource.LogInterface
}

func NewLogInterface(logInterfaceSource LogInterfaceSource.LogInterface) LogInterfaceRepository.LogInterface {
	return &LogInterfaceImpl{
		logInterfaceSource: logInterfaceSource,
	}
}

func (r *LogInterfaceImpl) CreateLog(logParam LogParamEntity.LogInterface) error {

	err := r.logInterfaceSource.Insert(logParam.ToModel())
	if err != nil {
		return err
	}

	return nil
}
