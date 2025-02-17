package cache_data_source

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	LoggerPackage "e-commerce/pkg/logger"
)

type Redis interface {
	GetConnection() *redis.Client
}

type RedisImpl struct {
	connection *redis.Client
	library    Library.Library
}

func New(
	config Config.Config,
	library Library.Library,
) Redis {
	path := "Redis:New"
	redisDb, err := library.Atoi(config.GetConfig().Redis.Database)
	if err != nil {
		err = CustomErrorPackage.New(Constants.ErrConfigurationRedis, err, path, library)
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:  path,
			Constants.Title: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}).Panic(err.(*CustomErrorPackage.CustomError).GetPlain())
	}

	connection := redis.NewClient(&redis.Options{
		Addr:     library.Sprintf("%s:%s", config.GetConfig().Redis.Host, config.GetConfig().Redis.Port),
		Password: config.GetConfig().Redis.Password,
		DB:       redisDb,
	})

	return &RedisImpl{
		connection: connection,
	}
}

func (r *RedisImpl) GetConnection() *redis.Client {
	return r.connection
}
