package main

import (
	"todo-list-api/internal/config"
	dbDriver "todo-list-api/internal/db"
	"todo-list-api/internal/logger"

	_ "github.com/lib/pq"
)

func main() {
	c, err := config.LoadConfig()
	logger.SetupLogger(c.ConfigYaml.Env)
	if err != nil {
		logger.Logger.Error("error getting congig", "error", err)
	}

	err = dbDriver.ApplyMigrations(c.ConfigYaml.Database)
	if err != nil {
		logger.Logger.Error("error aplying migration", "error", err)
	}
}
