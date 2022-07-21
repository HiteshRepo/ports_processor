package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/db"
	"github.com/hiteshpattanayak-tw/ports_processor/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

const (
	cutSet       = "file://"
	databaseName = "postgres"
)

type Migrator struct {
	pgMigrate *migrate.Migrate
	logger.Logger
}

func ProvideMigrator(config db.DatabaseConfig, db *gorm.DB, logger logger.Logger) (*Migrator, error) {
	dbConn, err := (db).DB()
	if err != nil {
		return nil, err
	}
	pgMigrate, err := initMigrate(dbConn, config.MigrationPath)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		pgMigrate: pgMigrate,
		Logger:    logger,
	}, nil
}

func (m Migrator) RunMigrations() {
	m.RunMigrationsWith(m.pgMigrate, "Postgres Database")
}

func (m Migrator) RunMigrationsWith(migrateInstance *migrate.Migrate, dBName string) {
	if err := migrateInstance.Up(); err != nil {
		if err == migrate.ErrNoChange {
			m.Logger.Debug(fmt.Sprintf("No change detected after running the migrations for %s", dBName))
			return
		}
		m.Logger.Fatal(fmt.Sprintf("Migration Failed for %s", dBName), zap.Error(err))
	}
	m.Logger.Info(fmt.Sprintf("Migrations applied successfully to %s", dBName))
}

func initMigrate(dbConn *sql.DB, directory string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	sourcePath, err := getSourcePath(directory)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(sourcePath, databaseName, driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func getSourcePath(directory string) (string, error) {
	directory = strings.TrimPrefix(directory, cutSet)

	absPath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", cutSet, absPath), nil
}
