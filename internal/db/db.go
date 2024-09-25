package db

import (
	"database/sql"
	"fmt"
	"todo-list-api/internal/config"
	"todo-list-api/internal/logger"

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
