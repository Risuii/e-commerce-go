package repository

import (
	LogEntity "e-commerce/internal/logging/domain/entity"
)

type LogActivity interface {
	CreateLog(logParam LogEntity.LogActivity) error
}
