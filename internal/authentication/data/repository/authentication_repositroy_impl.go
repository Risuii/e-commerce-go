package repository

import (
	"time"

	AuthenticationSource "e-commerce/internal/authentication/data/source"
	AuthenticationEntity "e-commerce/internal/authentication/domain/entity"
	AuthRepository "e-commerce/internal/authentication/domain/repository"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
)

type AuthenticationRepositoryImpl struct {
	library              Library.Library
	AuthenticationSource AuthenticationSource.AuthenticationMemory
}

func NewAuthenticationRepository(
	library Library.Library,
	AuthenticationSource AuthenticationSource.AuthenticationMemory,
) AuthRepository.AuthenticationRepository {
	return &AuthenticationRepositoryImpl{
		library:              library,
		AuthenticationSource: AuthenticationSource,
	}
}

func (r *AuthenticationRepositoryImpl) SetJWEToken(param *AuthenticationEntity.JWEToken, expiration time.Duration) error {
	path := "AuthenticationRepository:SetJWEToken"

	// INIT AUTHENTICATION MEMORY
	memory := r.AuthenticationSource
	// SET DATA TO CACHE
	err := memory.SetJWEToken(param.ToModel(), expiration)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	// SUCCESS SET DATA
	return nil
}

func (r *AuthenticationRepositoryImpl) GetJWETokenByKey(key string) (*AuthenticationEntity.JWEToken, error) {
	path := "AuthenticationRepository:GetJWETokenByKey"

	// INIT AUTHENTICATION MEMORY
	memory := r.AuthenticationSource

	// GET JWE TOKEN FROM REDIS
	model, err := memory.GetJWETokenByKey(key)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	// IF DATA NIL THEN RETURN NIL
	if model == nil {
		return nil, nil
	}

	// IF EXIST THEN RETURN TOKEN
	entity := AuthenticationEntity.JWEToken{}
	return entity.FromModel(model), nil
}

func (r *AuthenticationRepositoryImpl) DeleteJWETokenByKey(key string) error {
	path := "AuthenticationRepository:DeleteJWETokenByKey"

	// INIT AUTHENTICATION MEMORY
	memory := r.AuthenticationSource

	model, err := memory.GetJWETokenByKey(key)
	// ERROR HANDLER
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	if model == nil {
		return nil
	}

	err = memory.DeleteJWETokenByKey(key)
	// WHEN DELETE PROCESS RETURNS ERROR
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}
	// WHEN DELETE PROCESS IS SUCCESSFUL
	return nil
}
