package usecase

import (
	"net/http"
	"time"

	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/authentication/delivery/dto"
	StoreDTO "e-commerce/internal/store/delivery/dto"
	StoreEntity "e-commerce/internal/store/domain/entity"
	StoreRepository "e-commerce/internal/store/domain/repository"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
)

type UpdateStoreUsecase interface {
	Index(param *StoreDTO.StoreParam, credential *AuthDTO.LogoutParam) error
}

type UpdateStoreUsecaseImpl struct {
	library         Library.Library
	storeRepository StoreRepository.StoreRepository
}

func NewUpdateStoreUsecase(
	library Library.Library,
	storeRepository StoreRepository.StoreRepository,
) UpdateStoreUsecase {
	return &UpdateStoreUsecaseImpl{
		library:         library,
		storeRepository: storeRepository,
	}
}

func (u *UpdateStoreUsecaseImpl) Index(param *StoreDTO.StoreParam, credential *AuthDTO.LogoutParam) error {
	path := "CreateStoreUsecase:Index"

	// GET STORE DETAIL
	store, err := u.GetStore(credential.UserID)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// CHECKING STORE IS EXIST
	if store == nil {
		err := CustomErrorPackage.New(Constants.ErrStoreNotFound, Constants.ErrStoreNotFound, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusNotFound)
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// UPDATE STORE
	err = u.UpdateStore(store.StoreID, param.StoreName, param.Description)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (u *UpdateStoreUsecaseImpl) GetStore(userID string) (*StoreEntity.Store, error) {
	path := "CreateStoreUsecase:GetStore"

	// FIND STORE BY USER ID IN DB
	entity, err := u.storeRepository.GetStore(userID)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return entity, nil
}

func (u *UpdateStoreUsecaseImpl) UpdateStore(storeID, storeName, description string) error {
	path := "CreateStoreUsecase:UpdateStore"

	// SET ENTITY STORE
	entity := StoreEntity.Store{
		StoreID:     storeID,
		StoreName:   storeName,
		Description: description,
		UpdatedAt:   time.Now().Format(Constants.YYYMMDDHHMMSS),
	}

	// UPDATE DATA TO DB
	if err := u.storeRepository.UpdateStore(&entity); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}
