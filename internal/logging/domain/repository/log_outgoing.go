package repository

import (
	LogEntity "e-commerce/internal/logging/domain/entity"
)

type LogOutgoing interface {
	CreateLog(logParam LogEntity.LogOutgoing) error
}
