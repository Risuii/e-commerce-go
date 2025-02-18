package repository

import (
	UserEntity "e-commerce/internal/user/domain/entity"
)

type UserRepository interface {
	GetDetailUsers(username, email string) (*UserEntity.User, error)
	Insert(param *UserEntity.User) error
	UpdateLastLogin(userID, lastLogin string) error
}
