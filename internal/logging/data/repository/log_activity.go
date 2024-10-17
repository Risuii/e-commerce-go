package repository

import (
	LogActivitySource "e-commerce/internal/logging/data/source"
	LogParamEntity "e-commerce/internal/logging/domain/entity"
	LogActivityRepository "e-commerce/internal/logging/domain/repository"
)

type LogActivityImpl struct {
	logActivitySource LogActivitySource.LogActivity
}

func NewLogActivity(logAcitivitySource LogActivitySource.LogActivity) LogActivityRepository.LogActivity {
	return &LogActivityImpl{
		logActivitySource: logAcitivitySource,
	}
}

func (r *LogActivityImpl) CreateLog(logParam LogParamEntity.LogActivity) error {

	err := r.logActivitySource.Insert(logParam.ToModel())
	if err != nil {
		return err
	}

	return nil
}
