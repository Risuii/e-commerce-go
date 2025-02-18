package repository

import (
	UserSource "e-commerce/internal/user/data/source"
	UserEntity "e-commerce/internal/user/domain/entity"
	UserRepository "e-commerce/internal/user/domain/repository"
	CustomErrorPackage "e-commerce/pkg/custom_error"
)

type UserImpl struct {
	userSource UserSource.UserPersistent
}

func NewUserRepository(userSource UserSource.UserPersistent) UserRepository.UserRepository {
	return &UserImpl{
		userSource: userSource,
	}
}

func (r *UserImpl) GetDetailUsers(username, email string) (*UserEntity.User, error) {
	path := "UserRepository:GetDetailUsers"

	persistent := r.userSource
	model, err := persistent.GetDetail(username, email)
	if err != nil {
		return nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	if model == nil {
		return nil, nil
	}

	entity := UserEntity.User{}
	return entity.FromModel(model), nil
}

func (r *UserImpl) Insert(param *UserEntity.User) error {
	path := "UserRepository:Insert"

	persistent := r.userSource
	err := persistent.Create(param.ToModel())
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (r *UserImpl) UpdateLastLogin(userID, lastLogin string) error {
	path := "UserRepository:UpdateLastLogin"

	persistent := r.userSource
	err := persistent.UpdateLastLogin(userID, lastLogin)
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}
