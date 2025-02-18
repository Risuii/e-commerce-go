package repository

import (
	StoreEntity "e-commerce/internal/store/domain/entity"
)

type StoreRepository interface {
	CreateStore(param *StoreEntity.Store) error
	UpdateStore(param *StoreEntity.Store) error
	GetStore(userID string) (*StoreEntity.Store, error)
}
