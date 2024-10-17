package main

import (
	"context"
	"log"
	"os"

	Config "e-commerce/config"
	ecommerce "e-commerce/pkg/data_sources/e-commerce"

	"e-commerce/library"
	"e-commerce/migration"
)

func main() {
	// Initialize the context
	ctx := context.Background()

	// Initialize your library (if any dependency injection needed)
	lib := library.New()

	// Load your config (assuming you have some config loading mechanism)
	config := Config.New(lib)
	err := config.Setup()
	if err != nil {
		log.Println("Failed to load config")
	}

	// Initialize the Ecommerce database connection (singleton pattern)
	ecommerceDB := ecommerce.New(config, lib)

	// Get the GORM connection
	dbConnection := ecommerceDB.GetConnection()

	// Initialize the migration service with the GORM connection
	migrationService, err := migration.New(ctx, dbConnection)
	if err != nil {
		log.Fatalf("Failed to initialize migration service: %v", err)
	}

	args := os.Args
	if len(args) < 2 {
		log.Println("Missing args. args: [up | rollback]")
	}

	// Run migrations
	switch args[1] {
	case "up":
		migrationService.Up(ctx)
	case "rollback":
		migrationService.Rollback(ctx)
	default:
		log.Println("Invalid migration command")
	}

	log.Println("Migrations applied successfully!")
}
