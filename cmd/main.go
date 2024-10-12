package main

import (
	"net/http"
	"os"
	"todo-list-api/internal/config"
	"todo-list-api/internal/db"
	l "todo-list-api/internal/logger"
	"todo-list-api/internal/repository/postgres"
	"todo-list-api/internal/server"
	"todo-list-api/internal/server/handlers"
	"todo-list-api/internal/server/middleware"
	"todo-list-api/internal/service/task"
	"todo-list-api/internal/service/user"
)

func main() {
	c, err := config.LoadConfig()
	l.SetupLogger(c.ConfigYaml.Env)
	if err != nil {
		l.Logger.Error("error reading config file", "can't read config file", err)
		os.Exit(1)

	}
	l.Logger.Info("successfully loaded config", "env", c.ConfigYaml.Env)
	db, err := db.ConnectDB(c.ConfigYaml.Database)
	if err != nil {
		l.Logger.Error("error connecting to database", "error", err)
		os.Exit(1)
	}
	l.Logger.Info("successfully connected to database", "db_name", c.ConfigYaml.Database.Name)

	rep := postgres.NewPostgresRepo(db)
	taskService := task.NewTaskService(rep)
	userService := user.NewUserService(rep, c.ConfigEnv.SecretJWT)
	serverTask := handlers.NewTaskServer(*taskService)
	serverUser := handlers.NewUserServer(*userService)
	router := server.NewRouter(serverTask, serverUser)
	stack := middleware.CreateStack(
		middleware.LoggerMiddleware,
		middleware.TokenAuthMiddleware(c.ConfigEnv.SecretJWT),
	)
	server := &http.Server{
		Addr:         c.ConfigYaml.Server.Host + ":" + c.ConfigYaml.Server.Port,
		Handler:      stack(router),
		ReadTimeout:  c.ConfigYaml.Server.Timeout,
		WriteTimeout: c.ConfigYaml.Server.Timeout,
		IdleTimeout:  c.ConfigYaml.Server.IdleTimeout,
	}
	l.Logger.Info("connecting server", "address", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		l.Logger.Error("error running server", "error", err)
	}

}
