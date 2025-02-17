package repository

import (
	"time"

	AuthenticationEntity "e-commerce/internal/authentication/domain/entity"
)

type AuthenticationRepository interface {
	SetJWEToken(param *AuthenticationEntity.JWEToken, expiration time.Duration) error
	GetJWETokenByKey(key string) (*AuthenticationEntity.JWEToken, error)
	DeleteJWETokenByKey(key string) error
}
