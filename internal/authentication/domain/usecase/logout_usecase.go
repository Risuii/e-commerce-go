package usecase

import (
	Constants "e-commerce/constants"
	AuthDTO "e-commerce/internal/authentication/delivery/dto"
	AuthenticationRepository "e-commerce/internal/authentication/domain/repository"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	UtilsPackage "e-commerce/pkg/utils"
)

type LogoutUsecase interface {
	Index(param *AuthDTO.LogoutParam) error
}

type LogoutUsecaseImpl struct {
	library                  Library.Library
	authenticationRepository AuthenticationRepository.AuthenticationRepository
}

func NewLogoutUsecase(
	library Library.Library,
	authenticationRepository AuthenticationRepository.AuthenticationRepository,
) LogoutUsecase {
	return &LogoutUsecaseImpl{
		library:                  library,
		authenticationRepository: authenticationRepository,
	}
}

func (u *LogoutUsecaseImpl) Index(param *AuthDTO.LogoutParam) error {
	path := "LogoutUsecase:Index"

	// DELETE TOKEN IN REDIS
	if err := u.authenticationRepository.DeleteJWETokenByKey(Constants.UsernamePrefix + UtilsPackage.TernaryOperator(param.Username != Constants.NilString, param.Username, param.Email)); err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}
