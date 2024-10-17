package repository

import (
	LogOutgoingSource "e-commerce/internal/logging/data/source"
	LogParamEntity "e-commerce/internal/logging/domain/entity"
	LogOutgoingRepository "e-commerce/internal/logging/domain/repository"
)

type LogOutgoingImpl struct {
	logOutgoingSource LogOutgoingSource.LogOutgoing
}

func NewLogOutgoing(logOutgoingSource LogOutgoingSource.LogOutgoing) LogOutgoingRepository.LogOutgoing {
	return &LogOutgoingImpl{
		logOutgoingSource: logOutgoingSource,
	}
}

func (r *LogOutgoingImpl) CreateLog(logParam LogParamEntity.LogOutgoing) error {

	err := r.logOutgoingSource.Insert(logParam.ToModel())
	if err != nil {
		return err
	}

	return nil
}
