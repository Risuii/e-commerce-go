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

type ChangeStoreStatusUsecase interface {
	Index(param *StoreDTO.StoreStatusParam, credential *AuthDTO.LogoutParam) error
}

type ChangeStoreStatusUsecaseImpl struct {
	library         Library.Library
	storeRepository StoreRepository.StoreRepository
}

func NewChangeStoreStatus(
	library Library.Library,
	storeRepository StoreRepository.StoreRepository,
) ChangeStoreStatusUsecase {
	return &ChangeStoreStatusUsecaseImpl{
		library:         library,
		storeRepository: storeRepository,
	}
}

func (u *ChangeStoreStatusUsecaseImpl) Index(param *StoreDTO.StoreStatusParam, credential *AuthDTO.LogoutParam) error {
	path := "ChangeStoreStatusUsecase:Index"

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

	// UPDATE STORE STATUS
	err = u.UpdateStoreStatus(store.StoreID, store.StoreName, store.Description, param.Status)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (u *ChangeStoreStatusUsecaseImpl) GetStore(userID string) (*StoreEntity.Store, error) {
	path := "ChangeStoreStatusUsecase:GetStore"

	// FIND STORE BY USER ID IN DB
	entity, err := u.storeRepository.GetStore(userID)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return entity, nil
}

func (u *ChangeStoreStatusUsecaseImpl) UpdateStoreStatus(storeID, storeName, description, status string) error {
	path := "ChangeStoreStatusUsecase:UpdateStore"

	// SET ENTITY STORE
	entity := StoreEntity.Store{
		StoreID:     storeID,
		StoreName:   storeName,
		Description: description,
		Status:      status,
		UpdatedAt:   time.Now().Format(Constants.YYYMMDDHHMMSS),
	}

	// UPDATE DATA TO DB
	if err := u.storeRepository.UpdateStore(&entity); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}
