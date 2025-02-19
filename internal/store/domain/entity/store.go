package entity

import (
	StoreModel "e-commerce/internal/store/data/model"
)

type Store struct {
	StoreID     string
	StoreName   string
	Description string
	UserID      string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

func (e Store) FromModel(model *StoreModel.Store) *Store {
	result := Store{
		StoreID:     model.StoreID,
		StoreName:   model.StoreName,
		Description: model.Description,
		UserID:      model.UserID,
		Status:      model.Status,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}

	return &result
}

func (e Store) ToModel() *StoreModel.Store {
	result := StoreModel.Store{
		StoreID:     e.StoreID,
		StoreName:   e.StoreName,
		Description: e.Description,
		UserID:      e.UserID,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	return &result
}
