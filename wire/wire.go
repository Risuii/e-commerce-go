//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	Config "e-commerce/config"
	Library "e-commerce/library"
	BcryptPackage "e-commerce/pkg/bcrypt"
	CryptoPackage "e-commerce/pkg/crypto"
	CustomValidationPackage "e-commerce/pkg/custom_validation"
	JWEPackage "e-commerce/pkg/jwe"

	EcommercePackage "e-commerce/pkg/data_sources/e-commerce"
	RedisPackage "e-commerce/pkg/data_sources/redis"

	Routes "e-commerce/routes"

	AuthenticationRepository "e-commerce/internal/authentication/data/repository"
	AuthenticationSource "e-commerce/internal/authentication/data/source"
	AuthenticationHandler "e-commerce/internal/authentication/delivery/presenter/http"
	AuthenticationUsecase "e-commerce/internal/authentication/domain/usecase"

	UserRepository "e-commerce/internal/user/data/repository"
	UserSource "e-commerce/internal/user/data/source"

	StoreRepository "e-commerce/internal/store/data/repository"
	StoreSource "e-commerce/internal/store/data/source"
	StoreHandler "e-commerce/internal/store/delivery/presenter/http"
	StoreUsecase "e-commerce/internal/store/domain/usecase"

	ProductHandler "e-commerce/internal/product/delivery/presenter/http"

	LogUsecase "e-commerce/internal/logging/domain/usecase"

	ActivityLogRepository "e-commerce/internal/logging/data/repository"
	ActivityLogSource "e-commerce/internal/logging/data/source"

	Middleware "e-commerce/middlewares"
)

var ProviderSet = wire.NewSet(
	// FRAMEWORK
	gin.New,

	// PACKAGE
	CryptoPackage.NewCustomCrypto,
	BcryptPackage.NewBcrypt,
	CustomValidationPackage.NewCustomValidation,
	JWEPackage.NewJWE,

	// DATABASE
	EcommercePackage.New,
	RedisPackage.New,

	// DATASOURCE
	ActivityLogSource.NewLogActivityPersistent,
	UserSource.NewUserPersistent,
	AuthenticationSource.NewAuthenticationMemory,
	StoreSource.NewStorePersistent,

	// REPOSITORY
	ActivityLogRepository.NewLogActivity,
	UserRepository.NewUserRepository,
	AuthenticationRepository.NewAuthenticationRepository,
	StoreRepository.NewStoreRepository,

	// USECASE
	AuthenticationUsecase.NewRegisterUseCase,
	AuthenticationUsecase.NewLoginUsecase,
	AuthenticationUsecase.NewLogoutUsecase,
	LogUsecase.NewLogUsecase,
	StoreUsecase.NewCreateStoreUsecase,
	StoreUsecase.NewUpdateStoreUsecase,
	StoreUsecase.NewGetStoreUsecase,
	StoreUsecase.NewChangeStoreStatus,

	// HANDLER
	AuthenticationHandler.NewUserHandler,
	StoreHandler.NewStoreHandler,
	ProductHandler.NewProductHandler,

	// PUBLISHER

	// CONSUMER

	// MIDDLEWARE
	Middleware.NewMiddleware,

	// ROUTE
	Routes.New,

	// QUEUE
)

func InjectRoute(config Config.Config, library Library.Library) Routes.Routes {
	wire.Build(
		ProviderSet,
	)
	return nil
}
