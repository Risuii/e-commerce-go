package repository

import (
	UserSource "e-commerce/internal/user/data/source"
	UserEntity "e-commerce/internal/user/domain/entity"
	UserRepository "e-commerce/internal/user/domain/repository"
)

type UserImpl struct {
	userSource UserSource.User
}

func NewUser(userSource UserSource.User) UserRepository.UserRepository {
	return &UserImpl{
		userSource: userSource,
	}
}

func (r *UserImpl) GetDetailUsers(username, email string) (*UserEntity.User, error) {
	persistent := r.userSource
	userData, err := persistent.GetDetail(username, email)
	if err != nil {
		return nil, err
	}

	return (*UserEntity.User)(userData), nil
}

func (r *UserImpl) Insert(param *UserEntity.User) error {

	err := r.userSource.Create(param.ToModel())
	if err != nil {
		return err
	}

	return nil
}
