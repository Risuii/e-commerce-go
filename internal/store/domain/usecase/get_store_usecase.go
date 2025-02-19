package usecase

import (
	"net/http"

	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/authentication/delivery/dto"
	StoreDTO "e-commerce/internal/store/delivery/dto"
	StoreEntity "e-commerce/internal/store/domain/entity"
	StoreRepository "e-commerce/internal/store/domain/repository"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
)

type GetStoreUsecase interface {
	Index(credential *AuthDTO.LogoutParam) (*StoreDTO.StoreResponse, error)
}

type GetStoreUsecaseImpl struct {
	library         Library.Library
	storeRepository StoreRepository.StoreRepository
}

func NewGetStoreUsecase(
	library Library.Library,
	storeRepository StoreRepository.StoreRepository,
) GetStoreUsecase {
	return &GetStoreUsecaseImpl{
		library:         library,
		storeRepository: storeRepository,
	}
}

func (u *GetStoreUsecaseImpl) Index(credential *AuthDTO.LogoutParam) (*StoreDTO.StoreResponse, error) {
	path := "GetStoreUsecase:Index"

	// GET STORE DETAIL
	store, err := u.GetStore(credential.UserID)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// CHECKING STORE IS EXIST
	if store == nil {
		err := CustomErrorPackage.New(Constants.ErrStoreNotFound, Constants.ErrStoreNotFound, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusNotFound)
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// CHECKING STORE STATUS
	if store.Status == Constants.StoreStatusNotActive {
		err := CustomErrorPackage.New(Constants.ErrStoreNotFound, Constants.ErrStoreNotFound, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusNotFound)
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	response := StoreDTO.StoreResponse{
		StoreName:   store.StoreName,
		Description: store.Description,
	}

	return &response, nil
}

func (u *GetStoreUsecaseImpl) GetStore(userID string) (*StoreEntity.Store, error) {
	path := "GetStoreUsecase:GetStore"

	// FIND STORE BY USER ID IN DB
	entity, err := u.storeRepository.GetStore(userID)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return entity, nil
}
