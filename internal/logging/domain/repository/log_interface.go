package repository

import (
	LogEntity "e-commerce/internal/logging/domain/entity"
)

type LogInterface interface {
	CreateLog(logParam LogEntity.LogInterface) error
}
