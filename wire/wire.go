//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	Config "e-commerce/config"
	Library "e-commerce/library"
	CryptoPackage "e-commerce/pkg/crypto"
	Logging "e-commerce/pkg/logger"

	Routes "e-commerce/routes"

	ActivityLogRepository "e-commerce/internal/logging/data/repository"
	ActivityLogSource "e-commerce/internal/logging/data/source"

	InterfaceLogRepository "e-commerce/internal/logging/data/repository"
	InterfaceLogSource "e-commerce/internal/logging/data/source"

	OutgoingLogRepository "e-commerce/internal/logging/data/repository"
	OutgoingLogSource "e-commerce/internal/logging/data/source"
)

var ProviderSet = wire.NewSet(
	// FRAMEWORK
	gin.New,

	// PACKAGE
	CryptoPackage.NewCustomCrypto,
	Logging.NewLogger,

	// DATABASE

	// DATASOURCE
	ActivityLogSource.NewLogActivityPersistent,
	InterfaceLogSource.NewLogInterfacePersistent,
	OutgoingLogSource.NewLogOutgoingPersistent,

	// REPOSITORY
	OutgoingLogRepository.NewLogOutgoing,
	InterfaceLogRepository.NewLogInterface,
	ActivityLogRepository.NewLogActivity,

	// USECASE

	// HANDLER

	// PUBLISHER

	// CONSUMER

	// MIDDLEWARE

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
