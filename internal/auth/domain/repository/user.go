package repository

import (
	UserEntity "e-commerce/internal/auth/domain/entity"
)

type UserRepository interface {
	FindByUsername(param string) (*UserEntity.User, error)
	Insert(param *UserEntity.User) error
}
