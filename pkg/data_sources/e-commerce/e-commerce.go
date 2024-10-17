package ecommerce

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"

	GormPostgreServer "gorm.io/driver/postgres"
)

type Ecommerce interface {
	GetConnection() *gorm.DB
}

type EcommerceImpl struct {
	connection *gorm.DB
	config     Config.Config
	setup      string
	library    Library.Library
}

var (
	EcommerceInit EcommerceImpl
	EcommerceOnce sync.Once
)

func New(config Config.Config, library Library.Library) Ecommerce {
	EcommerceOnce.Do(func() {
		path := Constants.EcommerceDB
		// Setup configuration string
		setup := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			config.GetConfig().DBCon.Host,
			config.GetConfig().DBCon.User,
			config.GetConfig().DBCon.Password,
			config.GetConfig().DBCon.Database,
			config.GetConfig().DBCon.Port,
		)

		// Open connection
		connection, err := gorm.Open(GormPostgreServer.Open(setup), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		// Handle connection error
		if err != nil {
			err = CustomErrorPackage.New(Constants.ErrConnectionPostgres, err, path, library)
			log.Println(err)
			return
		}

		// Get underlying sql.DB object
		db, err := connection.DB()
		if err != nil {
			err = CustomErrorPackage.New(Constants.ErrConnectionPostgres, err, path, library)
			log.Println(err)
			return
		}

		// Configure database connection pool
		db.SetMaxIdleConns(config.GetConfig().DB.MaxIdleConns)
		db.SetMaxOpenConns(config.GetConfig().DB.MaxOpenConns)
		db.SetConnMaxLifetime(time.Duration(config.GetConfig().DB.ConnMaxLifetime) * time.Second)

		// Initialize singleton instance
		EcommerceInit = EcommerceImpl{
			connection: connection,
			config:     config,
			setup:      setup,
			library:    library,
		}
	})
	return &EcommerceInit
}

func (c *EcommerceImpl) GetConnection() *gorm.DB {
	return c.connection
}
