package migration

import (
	"context"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"

	Constants "e-commerce/constants"
)

type MigrationService interface {
	Up(context.Context) error
	Rollback(context.Context) error
	Version(context.Context) (int, bool, error)
}

type migrationService struct {
	migrate *migrate.Migrate
}

func New(ctx context.Context, db *gorm.DB) (MigrationService, error) {
	sqlDB, err := db.DB() // Get the *sql.DB from GORM
	if err != nil {
		log.Printf("error getting sql.DB from GORM: %s", err)
		return nil, err
	}

	databaseInstance, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Printf("go-migrate postgres drv init failed: %s", err)
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(Constants.FileMigration, Constants.MigrateLogIdentifier, databaseInstance)
	if err != nil {
		log.Printf("migrate init failed %s", err)
		return nil, err
	}

	return &migrationService{
		migrate: m,
	}, nil
}

func (s *migrationService) Up(ctx context.Context) error {
	if err := s.migrate.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("failed running migrations: %s", err)
		return err
	}
	log.Printf("Migration completed successfully.")
	return nil
}

func (s *migrationService) Rollback(ctx context.Context) error {
	if err := s.migrate.Steps(-1); err != nil {
		log.Printf("failed to rollback: %s", err)
		return err
	}
	log.Printf("Rollback completed successfully.")
	return nil
}

func (s *migrationService) Version(ctx context.Context) (int, bool, error) {
	version, dirty, err := s.migrate.Version()
	if err != nil {
		log.Printf("failed to get version: %s", err)
		return 0, false, err
	}
	return int(version), dirty, nil
}
