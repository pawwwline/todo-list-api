package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"todo-list-api/internal/config"
	"todo-list-api/internal/logger"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
)

func ConnectDB(dbCfg config.Database) (*sql.DB, error) {
	ConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SSLMode)
	db, err := sql.Open("postgres", ConnStr)
	logger.Logger.Debug("connection string", "ConnStr", ConnStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getDirection() (string, error) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			return "up", nil
		case "down":
			return "down", nil
		default:
			return "", errors.New("not right direction: try up or down")
		}
	}
	return "", nil
}

func ApplyMigrations(dbCfg config.Database) error {
	ConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SSLMode)
	m, err := migrate.New("file://internal/db/migrations", ConnStr)
	if err != nil {
		logger.Logger.Error("failed to create migrate instance", "error", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	dir, err := getDirection()
	if err != nil {
		return err
	}
	if dir == "up" {
		err := m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
	} else {
		err := m.Down()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
	}
	version, dirty, _ := m.Version()
	logger.Logger.Info("Applied migration", "version", version, "dirty", dirty)
	return nil
}
