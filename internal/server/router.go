package server

import (
	"net/http"
	// l "todo-list-api/internal/logger"
)

type Router struct {
	taskServer TaskServerInterface
	userServer UserServerInterface
	mux        *http.ServeMux
}

func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewRouter(ts TaskServerInterface, us UserServerInterface) *Router {
	s := &Router{
		taskServer: ts,
		userServer: us,
		mux:        http.NewServeMux(),
	}
	s.initializeRoutes()
	return s
}

func (s *Router) initializeRoutes() {
	s.mux.HandleFunc("POST /register", s.userServer.SignupHandler)
	s.mux.HandleFunc("POST /login", s.userServer.LoginHandler)
	s.mux.HandleFunc("POST /todos", s.taskServer.CreateTaskHandler)
	s.mux.HandleFunc("PUT /todos/{id}", s.taskServer.UpdateTaskHandler)
	s.mux.HandleFunc("DELETE /todos/{id}", s.taskServer.DeleteTaskHandler)
	s.mux.HandleFunc("GET /todos", s.taskServer.GetAllTasksHandler)
}

// func (s *Server) Start(port string) {

// 	l.Logger.Info("starting server", "server port", port)
// 	err := http.ListenAndServe(":"+port, s.mux)
// 	if err != nil {
// 		l.Logger.Error("error starting server", "error", err)
// 	}
// }
