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

	EcommercePackage "e-commerce/pkg/data_sources/e-commerce"

	Routes "e-commerce/routes"

	UserRepository "e-commerce/internal/auth/data/repository"
	UserSource "e-commerce/internal/auth/data/source"
	UserHandler "e-commerce/internal/auth/delivery/presenter/http"
	UsertUseCase "e-commerce/internal/auth/domain/usecase"

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

	// DATABASE
	EcommercePackage.New,

	// DATASOURCE
	ActivityLogSource.NewLogActivityPersistent,
	UserSource.NewUserImpl,

	// REPOSITORY
	ActivityLogRepository.NewLogActivity,
	UserRepository.NewUser,

	// USECASE
	UsertUseCase.NewRegisterUseCase,
	LogUsecase.NewLogUsecase,

	// HANDLER
	UserHandler.NewUserHandler,

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
