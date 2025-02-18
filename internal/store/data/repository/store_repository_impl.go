package repository

import (
	StoreSource "e-commerce/internal/store/data/source"
	StoreEntity "e-commerce/internal/store/domain/entity"
	StoreRepository "e-commerce/internal/store/domain/repository"
	CustomErrorPackage "e-commerce/pkg/custom_error"
)

type StoreImpl struct {
	storeSource StoreSource.StorePersistent
}

func NewStoreRepository(
	storeSource StoreSource.StorePersistent,
) StoreRepository.StoreRepository {
	return &StoreImpl{
		storeSource: storeSource,
	}
}

func (r *StoreImpl) CreateStore(param *StoreEntity.Store) error {
	path := "StoreRepository:CreateStore"

	persistent := r.storeSource
	err := persistent.Create(param.ToModel())
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (r *StoreImpl) UpdateStore(param *StoreEntity.Store) error {
	path := "StoreRepository:UpdateStore"

	persistent := r.storeSource
	err := persistent.Update(param.ToModel())
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (r *StoreImpl) GetStore(userID string) (*StoreEntity.Store, error) {
	path := "StoreRepository:GetStore"

	persistent := r.storeSource
	model, err := persistent.GetStore(userID)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	if model == nil {
		return nil, nil
	}

	entity := StoreEntity.Store{}
	return entity.FromModel(model), nil
}
