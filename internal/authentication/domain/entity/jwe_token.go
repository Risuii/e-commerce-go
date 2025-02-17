package entity

import (
	AuthModel "e-commerce/internal/authentication/data/model"
)

type JWEToken struct {
	Username         string
	TOKEN            string
	JWT_TOKEN_PREFIX string
}

func (e JWEToken) FromModel(model *AuthModel.JWEToken) *JWEToken {
	e.TOKEN = model.TOKEN
	return &e
}

func (e JWEToken) ToModel() *AuthModel.JWEToken {
	result := AuthModel.JWEToken{
		KEY:   e.JWT_TOKEN_PREFIX + e.Username,
		TOKEN: e.TOKEN,
	}
	return &result
}
