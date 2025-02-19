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
	ExecutionResultPackage "e-commerce/pkg/execution_result"
)

type CreateStoreUsecase interface {
	Index(param *StoreDTO.StoreParam, credential *AuthDTO.LogoutParam) error
}

type CreateStoreUsecaseImpl struct {
	library         Library.Library
	storeRepository StoreRepository.StoreRepository
}

func NewCreateStoreUsecase(
	library Library.Library,
	storeRepository StoreRepository.StoreRepository,
) CreateStoreUsecase {
	return &CreateStoreUsecaseImpl{
		library:         library,
		storeRepository: storeRepository,
	}
}

func (u *CreateStoreUsecaseImpl) Index(param *StoreDTO.StoreParam, credential *AuthDTO.LogoutParam) error {
	path := "CreateStoreUsecase:Index"

	// CHECKING STORE IS EXIST
	storeChannel := make(chan ExecutionResultPackage.ExecutionResult)
	idChannel := make(chan ExecutionResultPackage.ExecutionResult)

	// DO GOROUTINE
	go u.FindStore(credential.UserID, storeChannel)
	go u.GenerateID(idChannel)

	// GET STORE DATA FROM CHANNEL
	storeDetail := <-storeChannel
	if err := storeDetail.GetError(); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	store := storeDetail.GetData().(*StoreEntity.Store)

	// CHECKING STORE IS EXIST
	if store != nil {
		err := CustomErrorPackage.New(u.MappingStoreStatus(store.Status), u.MappingStoreStatus(store.Status), path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusConflict)
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// GET ID DATA FROM CHANNEL
	idDetail := <-idChannel
	if err := idDetail.GetError(); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	storeID := idDetail.GetData().(string)

	// CHECKING STORE ID EXIST
	if storeID == Constants.NilString {
		err := CustomErrorPackage.New(Constants.ErrSomethingWentWrong, Constants.ErrSomethingWentWrong, path, u.library)
		err.(*CustomErrorPackage.CustomError).SetCode(http.StatusInternalServerError)
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// INSERT STORE DATA TO DB
	if err := u.InsertStore(storeID, param.StoreName, param.Description, credential.UserID); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (u *CreateStoreUsecaseImpl) GenerateID(resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "CreateStoreUsecase:GenerateID"

	result := ExecutionResultPackage.ExecutionResult{}

	// GENERATE UUID FOR ID
	id, err := u.library.GenerateUUID()
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// SET UUID TO CHANNEL
	result.SetResult(id, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *CreateStoreUsecaseImpl) FindStore(userID string, resultChannel chan ExecutionResultPackage.ExecutionResult) {
	path := "CreateStoreUsecase:FindStore"

	result := ExecutionResultPackage.ExecutionResult{}

	// FIND STORE BY USER ID IN DB
	entity, err := u.storeRepository.GetStore(userID)
	if err != nil {
		result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path))
		resultChannel <- result
		close(resultChannel)
		return
	}

	// SET DATA ENTITY TO CHANNEL
	result.SetResult(entity, nil)
	resultChannel <- result
	close(resultChannel)
}

func (u *CreateStoreUsecaseImpl) InsertStore(storeID, storeName, description, userID string) error {
	path := "CreateStoreUsecase:FindStore"

	// INIT ENTITY
	entity := StoreEntity.Store{
		StoreID:     storeID,
		StoreName:   storeName,
		Description: description,
		UserID:      userID,
		Status:      Constants.StoreStatusActive,
		CreatedAt:   time.Now().Format(Constants.YYYMMDDHHMMSS),
	}

	// INSERT ENTITY STORE TO DB
	err := u.storeRepository.CreateStore(&entity)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (u *CreateStoreUsecaseImpl) MappingStoreStatus(status string) error {
	err := map[string]error{
		Constants.StoreStatusActive:    Constants.ErrOneStore,
		Constants.StoreStatusNotActive: Constants.ErrInactiveStore,
	}

	return err[status]
}
