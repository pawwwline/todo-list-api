package server

import (
	"net/http"
)

type TaskServerInterface interface {
	CreateTaskHandler(w http.ResponseWriter, r *http.Request)
	DeleteTaskHandler(w http.ResponseWriter, r *http.Request)
	UpdateTaskHandler(w http.ResponseWriter, r *http.Request)
	GetAllTasksHandler(w http.ResponseWriter, r *http.Request)
}

type UserServerInterface interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	SignupHandler(w http.ResponseWriter, r *http.Request)
}
