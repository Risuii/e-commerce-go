package entity

import (
	UserModel "e-commerce/internal/auth/data/model"
)

type User struct {
	Uuid          string
	Email         string
	EmailVerified bool
	Password      string
	Username      string
	LastLogin     string
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
}

func (e User) FromModel(model *UserModel.User) *User {
	result := User{
		Uuid:          model.Uuid,
		Email:         model.Email,
		EmailVerified: model.EmailVerified,
		Password:      model.Password,
		Username:      model.Username,
		LastLogin:     model.LastLogin,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		DeletedAt:     model.DeletedAt,
	}

	return &result
}

func (e User) ToModel() *UserModel.User {
	result := UserModel.User{
		Uuid:          e.Uuid,
		Email:         e.Email,
		EmailVerified: e.EmailVerified,
		Password:      e.Password,
		Username:      e.Username,
		LastLogin:     e.LastLogin,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		DeletedAt:     e.DeletedAt,
	}

	return &result
}
