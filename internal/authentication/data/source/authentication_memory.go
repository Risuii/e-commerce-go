package source

import (
	"time"

	"github.com/go-redis/redis"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	AuthenticationModel "e-commerce/internal/authentication/data/model"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	RedisPackage "e-commerce/pkg/data_sources/redis"
)

type AuthenticationMemory interface {
	SetJWEToken(param *AuthenticationModel.JWEToken, expiration time.Duration) error
	GetJWETokenByKey(key string) (*AuthenticationModel.JWEToken, error)
	DeleteJWETokenByKey(key string) error
}

type AuthenticationMemoryImpl struct {
	config  Config.Config
	library Library.Library
	cache   RedisPackage.Redis
}

func NewAuthenticationMemory(
	config Config.Config,
	library Library.Library,
	cache RedisPackage.Redis,
) AuthenticationMemory {
	return &AuthenticationMemoryImpl{
		config:  config,
		library: library,
		cache:   cache,
	}
}

func (s *AuthenticationMemoryImpl) SetJWEToken(param *AuthenticationModel.JWEToken, expiration time.Duration) error {
	path := "AuthenticationMemory:SetJWEToken"

	// SET TOKEN TO REDIS
	err := s.cache.GetConnection().Set(param.KEY, param.TOKEN, expiration).Err()
	if err != nil {
		return CustomErrorPackage.New(Constants.ErrSomethingWentWrong, err, path, s.library)
	}

	// RETURN NIL IF SUCCESS
	return nil
}

func (s *AuthenticationMemoryImpl) GetJWETokenByKey(key string) (*AuthenticationModel.JWEToken, error) {
	path := "AuthenticationMemory:GetJWETokenByKey"
	// GET CACHE BY KEY
	data, err := s.cache.GetConnection().Get(key).Result()
	// IF NOT ERROR
	if err == nil {
		// SUCCESS GET CACHE
		model := AuthenticationModel.JWEToken{
			KEY:   key,
			TOKEN: data,
		}
		return &model, nil
	}
	// HANDLE ERROR REDIS NIL
	if err == redis.Nil {
		return nil, nil
	}
	// ERROR HANDLER
	return nil, CustomErrorPackage.New(Constants.ErrSomethingWentWrong, err, path, s.library)
}

func (s *AuthenticationMemoryImpl) DeleteJWETokenByKey(key string) error {
	path := "AuthenticationMemory:DeleteJWETokenByKey"

	// GET CACHE BY KEY
	deleted, err := s.cache.GetConnection().Del(key).Result()
	if err != nil {
		return CustomErrorPackage.New(Constants.ErrSomethingWentWrong, err, path, s.library)
	}
	// WHEN USER LOCK IS SUCCESSFULLY DELETED
	if deleted == 1 {
		return nil
	}
	// ERROR
	return CustomErrorPackage.New(Constants.ErrSomethingWentWrong, Constants.ErrFailedRemove, path, s.library)
}
